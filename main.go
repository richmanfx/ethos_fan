/***********************************************/
/* Программа для автоматического регулирования */
/* оборотов вентиляторов видеокарт в EthOS     */
/***********************************************/

package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/ini.v1"
	. "os"
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
		//gpuQuantity		int
	)

	// Выставить параметры логирования
	SetLog(log.DebugLevel)

	// Полное имя конфигурационного файла
	fullConfigFileName := configDirName + "/" + configFileName
	log.Debugf("Full config file name: %s", fullConfigFileName)

	// Чтение параметров из конфигурационного файла
	getConfigParameters(fullConfigFileName, debugLevel, highTemp, lowTemp, sleepTime, speedStep, initFanSpeed, minFanSpeed)

	// Получить количество GPU в системе

	// Выставить начальные скорости вентиляторов

	// Основной цикл

}

func getConfigParameters(
	fullConfigFileName string, debugLevel string,
	highTemp int, lowTemp int, sleepTime int, speedStep int, initFanSpeed int, minFanSpeed int) {

	config, err := ini.Load(fullConfigFileName)
	if err != nil {
		log.Debugf("Fail to read config file '%s': %v", fullConfigFileName, err)
		Exit(1)
	}

	debugLevel = config.Section("").Key("DEBUG_LEVEL").String()
	if debugLevel == "INFO" {
		SetLog(log.InfoLevel)
	}
	log.Debugf("Debug level: %s", debugLevel)

	highTemp = config.Section("").Key("HIGH_TEMP").MustInt(60)
	log.Debugf("High temperature: %d°C", highTemp)

	lowTemp = config.Section("").Key("LOW_TEMP").MustInt(55)
	log.Debugf("Low temperature: %d°C", lowTemp)

	sleepTime = config.Section("").Key("SLEEP_TIME").MustInt(60)
	log.Debugf("Sleep time: %ds", sleepTime)

	speedStep = config.Section("").Key("SPEED_STEP").MustInt(5)
	log.Debugf("Speed step: %d%%", speedStep)

	initFanSpeed = config.Section("").Key("INIT_FAN_SPEED").MustInt(80)
	log.Debugf("Initial fan Speed: %d%%", initFanSpeed)

	minFanSpeed = config.Section("").Key("MIN_FAN_SPEED").MustInt(15)
	log.Debugf("Minimal fan Speed: %d%%", minFanSpeed)
}

func SetLog(debugLevel log.Level) {
	log.SetOutput(Stdout)
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006/01/02 15:04:05"
	log.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
	log.SetLevel(debugLevel) // Уровень логирования, до уточнения из конфиг файла
}
