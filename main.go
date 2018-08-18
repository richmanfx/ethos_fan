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

	log.Debugf("GPUs: %d", gpuQuantity)

	log.Debugf("Temperature GPU0: %d", getGpuTemp(0))

	log.Debugf("Fan speed GPU0: %d", getGpuFanSpeed(0))

	//setGpuFanSpeed(4, 40)
	setInitialFanSpeed(gpuQuantity, initFanSpeed)

	// Выставить начальные скорости вентиляторов

	// Основной цикл

}

/* Первоначальная установка оборотов вентиляторов после старта системы */
func setInitialFanSpeed(gpuQuantity, initFanSpeed int) {

	for gpu := 0; gpu <= gpuQuantity-1; gpu++ {
		log.Debugf("GPU: %d, initFanSpeed: %d", gpu, initFanSpeed)
		setGpuFanSpeed(gpu, initFanSpeed)
	}

}

/* Выставить обороты вентиляторов GPU */
func setGpuFanSpeed(gpuNumber, newFanSpeed int) {
	// Нумерация GPU с нуля
	command := fmt.Sprintf("sudo ethos-smi --gpu %d --fan %d | cut -c 3-", gpuNumber, newFanSpeed)
	out, err := exec.Command("bash", "-c", command).Output()

	if err != nil {
		log.Debugf("Failed to execute command: %s", out)
	}

	log.Debugf("New GPU fan speed: '%s'", strings.Trim(string(out), "\n"))
}

// TODO: почти одинаковые функции - вынести в отдельную
/* Обороты вентилятора GPU */
func getGpuFanSpeed(gpuNumber int) (gpuFanSpeed int) {
	// Нумерация GPU с нуля
	command := fmt.Sprintf("ethos-smi -g %d | grep \"* Fan Speed\" | cut -f 5 -d \" \" | rev | cut -c 2- | rev", gpuNumber)
	out, err := exec.Command("bash", "-c", command).Output()

	if err != nil {
		log.Debugf("Failed to execute command: %s", out)
	}

	gpuFanSpeed, _ = strconv.Atoi(strings.Trim(string(out), "\n"))
	log.Debugf("GPU fan speed: '%d'", gpuFanSpeed)

	return
}

// TODO: почти одинаковые функции - вынести в отдельную
/* Температура GPU */
func getGpuTemp(gpuNumber int) (gpuTemp int) {
	// Нумерация GPU с нуля

	command := fmt.Sprintf("ethos-smi -g %d | grep \"* Temperature\" | cut -f 4 -d \" \" | rev | cut -c 2- | rev", gpuNumber)
	out, err := exec.Command("bash", "-c", command).Output()

	if err != nil {
		log.Debugf("Failed to execute command: %s", out)
	}

	gpuTemp, _ = strconv.Atoi(strings.Trim(string(out), "\n"))
	log.Debugf("GPU temperature: '%d'", gpuTemp)

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

	command := "/opt/ethos/bin/ethos-smi | grep \"\\[\" | grep \"\\]\" | grep GPU | tail -1 | cut -f 1 -d \" \" | cut -c 4,5"
	out, err := exec.Command("bash", "-c", command).Output()

	if err != nil {
		log.Debugf("Failed to execute command: %s", out)
	}

	gpuQuantity, _ = strconv.Atoi(strings.Trim(string(out), "\n"))

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
