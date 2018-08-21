/***********************************************/
/* Программа для автоматического регулирования */
/* оборотов вентиляторов видеокарт в EthOS     */
/***********************************************/

package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/ini.v1"
	. "os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {

	const (
		configDirName  = "/usr/local/etc/ethos_fan"
		configFileName = "fan.cfg"
	)

	var (
		debugLevel   = "DEBUG"
		highTemp     int
		lowTemp      int
		sleepTime    int
		speedStep    int
		initFanSpeed int
		minFanSpeed  int
		gpuQuantity  int
	)

	// Выставить параметры логирования
	SetLog(log.DebugLevel)

	// Полное имя конфигурационного файла
	fullConfigFileName := configDirName + "/" + configFileName
	log.Debugf("Full config file name: %s", fullConfigFileName)

	// Чтение параметров из конфигурационного файла
	getConfigParameters(
		fullConfigFileName, debugLevel, &highTemp, &lowTemp, &sleepTime, &speedStep, &initFanSpeed, &minFanSpeed)

	// Получить количество GPU в системе
	gpuQuantity = getGpuQuantity()

	// Выставить начальные скорости вентиляторов
	setInitialFanSpeed(gpuQuantity, initFanSpeed)

	// Основной цикл
	for {
		log.Debugln("-=========================================================-")
		log.Debugf("Cycle time: %d seconds.", sleepTime)
		setNewFanSpeedForAllGpu(gpuQuantity, initFanSpeed, lowTemp, highTemp, speedStep, minFanSpeed)

		time.Sleep(time.Duration(sleepTime) * time.Second)

		// TODO: Здесь нужно перехватывать SIGINT и выходить?
	}

}

/* Изменить скорость вентиляторов для всех карт в зависимости от температуры */
func setNewFanSpeedForAllGpu(gpuQuantity, initFanSpeed, lowTemp, highTemp, speedStep, minFanSpeed int) {

	var newFanSpeed int
	var currentTemp int
	var currentFanSpeed int

	newFanSpeed = initFanSpeed

	for gpu := 0; gpu <= gpuQuantity-1; gpu++ {

		currentTemp = getGpuTemp(gpu)
		currentFanSpeed = getGpuFanSpeed(gpu)
		log.Debugf("GPU %d: Temp: %d°C, Fan speed: %d%%", gpu, currentTemp, currentFanSpeed)

		if !checkValidInRange(lowTemp, highTemp, currentTemp) {

			log.Debugf("GPU %d: Out of temperature range %d...%d °C", gpu, lowTemp, highTemp)

			// Новая скорость кулеров при высокой температуре
			if currentTemp > highTemp {
				newFanSpeed = checkHighTemp(currentTemp, highTemp, currentFanSpeed, speedStep)
			}

			// Новая скорость кулеров при низкой температуре
			if currentTemp < lowTemp {
				newFanSpeed = checkLowTemp(currentTemp, lowTemp, currentFanSpeed, minFanSpeed, speedStep)
			}

			// Выставляем новую скорость
			if newFanSpeed != currentFanSpeed {
				//if (currentTemp > highTemp) || (currentTemp < lowTemp) {
				go setGpuFanSpeed(gpu, newFanSpeed)
			}

		}
	}
}

/* Новая скорость кулеров при температуре ниже нижнего уровня */
func checkLowTemp(currentTemp int, lowTemp int, currentFanSpeed int, minFanSpeed int, speedStep int) (newFanSpeed int) {

	if currentTemp < lowTemp {
		if currentFanSpeed > minFanSpeed {
			newFanSpeed = currentFanSpeed - speedStep
		} else {
			newFanSpeed = minFanSpeed
		}
	} else {
		newFanSpeed = currentFanSpeed
	}
	return newFanSpeed
}

/* Новая скорость кулеров при температуре выше верхнего уровня */
func checkHighTemp(currentTemp int, highTemp int, currentFanSpeed int, speedStep int) (newFanSpeed int) {

	if currentTemp > highTemp {
		if currentFanSpeed < 100 {
			if currentFanSpeed > 94 {
				newFanSpeed = 100
			} else {
				newFanSpeed = currentFanSpeed + speedStep
			}
		} else {
			newFanSpeed = 100
		}
	} else {
		newFanSpeed = currentFanSpeed
	}

	return newFanSpeed
}

/* Первоначальная установка оборотов вентиляторов после старта системы */
func setInitialFanSpeed(gpuQuantity, initFanSpeed int) {
	for gpu := 0; gpu <= gpuQuantity-1; gpu++ {
		log.Debugf("GPU %d: The initial speed of the fans: %d%%", gpu, initFanSpeed)
		setGpuFanSpeed(gpu, initFanSpeed)
	}
}

/* Выставить обороты вентиляторов GPU */
func setGpuFanSpeed(gpuNumber, newFanSpeed int) {
	// Нумерация GPU с нуля

	if newFanSpeed < 5 {
		log.Debugf("Alarm!!! New fan speed extremely low: '%v'", newFanSpeed)
	} else {
		command := fmt.Sprintf("sudo ethos-smi --gpu %d --fan %d | cut -c 3-", gpuNumber, newFanSpeed)
		out, err := exec.Command("bash", "-c", command).Output()

		if err != nil {
			log.Debugf("Failed to execute command: %s", out)
		}

		// Но всегда возвращается информация из ОС после изменения скорости (например для nVidea карт???)
		if len(strings.Trim(string(out), "\n")) < 5 {
			log.Debugf("GPU %d: Set %d%% fan.", gpuNumber, newFanSpeed)
		} else {
			log.Debugf("GPU %d: %s", gpuNumber, strings.Trim(string(out), "\n"))
		}

	}
}

/* Выполнить команду в ОС, вернуть результат */
func runOsCommand(commandParameter int, osCommandTemplate string) (result int) {
	command := fmt.Sprintf(osCommandTemplate, commandParameter)
	out, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		log.Debugf("Failed to execute command: %s", out)
	}
	result, _ = strconv.Atoi(strings.Trim(string(out), "\n"))
	return
}

/* Обороты вентилятора GPU */
func getGpuFanSpeed(gpuNumber int) (gpuFanSpeed int) {
	osCommandTemplate := "ethos-smi -g %d | grep \"* Fan Speed\" | cut -f 5 -d \" \" | rev | cut -c 2- | rev"
	gpuFanSpeed = runOsCommand(gpuNumber, osCommandTemplate)
	return
}

/* Температура GPU */
func getGpuTemp(gpuNumber int) (gpuTemp int) {
	osCommandTemplate := "ethos-smi -g %d | grep \"* Temperature\" | cut -f 4 -d \" \" | rev | cut -c 2- | rev"
	gpuTemp = runOsCommand(gpuNumber, osCommandTemplate)
	return
}

/*****************************************************************************
 * Проверяет попадание значения в валидный диапазон значений                 *
 *    Param: minimum - нижнее значение диапазона                             *
 *    Param: maximum - верхнее значение диапазона                            *
 *    Param: value - проверяемое значение                                 *
 *    Return: True - значение в диапазоне, False - значение вне диапазона    *
 *****************************************************************************/
func checkValidInRange(minimum, maximum, value int) (result bool) {
	if (value >= minimum) && (value <= maximum) {
		result = true
	} else {
		result = false
	}
	return result
}

/* Количество GPU в системе */
func getGpuQuantity() (gpuQuantity int) {
	osCommandTemplate := "/opt/ethos/bin/ethos-smi | grep \"\\[\" | grep \"\\]\" | grep GPU | tail -1 | cut -f %d -d \" \" | cut -c 4,5"
	gpuQuantity = runOsCommand(1, osCommandTemplate)
	gpuQuantity += 1 // Нумерация GPU в системе начинается с нуля
	log.Debugf("GPU quantity: '%d'", gpuQuantity)
	return
}

/* Получить параметры из конфигурационного INI файла */
func getConfigParameters(
	fullConfigFileName string, debugLevel string,
	highTemp *int, lowTemp *int, sleepTime *int, speedStep *int, initFanSpeed *int, minFanSpeed *int) {

	config, err := ini.Load(fullConfigFileName)
	if err != nil {
		log.Debugf("Fail to read config file: %v", err)
		Exit(1)
	}

	debugLevel = config.Section("").Key("DEBUG_LEVEL").String()
	if debugLevel == "INFO" {
		SetLog(log.InfoLevel)
	}
	log.Debugf("Debug level: %s", debugLevel)

	*highTemp = config.Section("").Key("HIGH_TEMP").MustInt(60)
	log.Debugf("High temperature: %d°C", *highTemp)

	*lowTemp = config.Section("").Key("LOW_TEMP").MustInt(55)
	log.Debugf("Low temperature: %d°C", *lowTemp)

	*sleepTime = config.Section("").Key("SLEEP_TIME").MustInt(60)
	log.Debugf("Sleep time: %ds", *sleepTime)

	*speedStep = config.Section("").Key("SPEED_STEP").MustInt(5)
	log.Debugf("Speed step: %d%%", *speedStep)

	*initFanSpeed = config.Section("").Key("INIT_FAN_SPEED").MustInt(80)
	log.Debugf("Initial fan Speed: %d%%", *initFanSpeed)

	*minFanSpeed = config.Section("").Key("MIN_FAN_SPEED").MustInt(15)
	log.Debugf("Minimal fan Speed: %d%%", *minFanSpeed)
}

/* Выставить параметры логирования */
func SetLog(debugLevel log.Level) {
	log.SetOutput(Stdout)
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006/01/02 15:04:05"
	log.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
	log.SetLevel(debugLevel)
}
