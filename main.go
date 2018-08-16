/***********************************************/
/* Программа для автоматического регулирования */
/* оборотов вентиляторов видеокарт в EthOS     */
/***********************************************/

package main

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/ini.v1"
	"io"
	. "os"
	"os/exec"
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
		fullConfigFileName, debugLevel,
		highTemp, lowTemp, sleepTime, speedStep, initFanSpeed, minFanSpeed)

	// Получить количество GPU в системе
	getGpuQuantity(gpuQuantity)

	// Выставить начальные скорости вентиляторов

	// Основной цикл

}

/* Количество GPU в системе */
func getGpuQuantity(gpuQuantity int) {
	gpuQuantity = 0

	//command := exec.Command("/opt/ethos/bin/ethos-smi | grep \"\\[\" | grep \"\\]\" | grep GPU | tail -1 | cut -f 1 -d \" \" | cut -c 4,5")
	//command := exec.Command("ls", "-la")
	command1 := exec.Command("ethos-smi")
	command2 := exec.Command("grep", "GPU")

	r, w := io.Pipe()
	command1.Stdout = w
	command2.Stdin = r

	var b2 bytes.Buffer
	command2.Stdout = &b2

	command1.Start()
	command2.Start()
	command1.Wait()
	w.Close()
	command2.Wait()
	io.Copy(Stdout, &b2)

	// Про пайп; https://stackoverflow.com/questions/10781516/how-to-pipe-several-commands-in-go
	//command := exec.Command("ethos-smi", "|", "grep", "GPU")
	//var buf bytes.Buffer
	//command.Stdout = &buf
	//err := command.Start()
	//if err != nil {
	//	log.Debugf("error: %v\n", err)
	//}
	//err = command.Wait()
	log.Debugf("Command finished with output: %v\n", b2.String())
	//log.Debugf("GPU quantity: %d", gpuQuantity)
}

/* Получить параметры из конфигурационного INI файла */
func getConfigParameters(
	fullConfigFileName string, debugLevel string,
	highTemp int, lowTemp int, sleepTime int, speedStep int, initFanSpeed int, minFanSpeed int) {

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

/* Выставить параметры логирования */
func SetLog(debugLevel log.Level) {
	log.SetOutput(Stdout)
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006/01/02 15:04:05"
	log.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
	log.SetLevel(debugLevel)
}
