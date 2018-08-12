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
	//highTemp int
	//lowTemp int
	//sleepTime int
	//speedStep int
	//initFanSpeed int
	//minFanSpeed int
	//gpuQuantity int
	)

	// Выставить параметры логирования
	SetLog(log.DebugLevel)

	// Полное имя конфигурационного файла
	fullConfigFileName := configDirName + "/" + configFileName
	log.Debugf("Full config file name: %s", fullConfigFileName)

	// Чтение параметров из конфигурационного файла
	config, err := ini.Load(fullConfigFileName)
	if err != nil {
		log.Debugf("Fail to read config file '%s': %v", fullConfigFileName, err)
		Exit(1)
	}
	log.Debugln("Debug level:", config.Section("").Key("DEBUG_LEVEL").String())

	// Получить количество GPU в системе

	// Выставить начальные скорости вентиляторов

	// Основной цикл

}

func SetLog(debugLevel log.Level) {
	log.SetOutput(Stdout)
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006/01/02 15:04:05"
	log.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
	log.SetLevel(debugLevel) // Уровень логирования, до уточнения из конфиг файла
}
