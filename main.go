/***********************************************/
/* Программа для автоматического регулирования */
/* оборотов вентиляторов видеокарт в EthOS     */
/***********************************************/

package main

import (
	log "github.com/Sirupsen/logrus"
)

func main() {

	// TODO: Вынести в конфиг
	log.SetLevel(log.DebugLevel) // Уровень логирования

	// Полное имя конфигурационного файла
	log.Debugln("Полное имя конфигурационного файла")

	// Чтение параметров из конфигурационного файла

	// Получить количество GPU в системе

	// Выставить начальные скорости вентиляторов

	// Основной цикл

}
