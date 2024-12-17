package pattern

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern
*/

type CPU struct {
}

func (c *CPU) Freeze()  {}
func (c *CPU) Execute() {}

type HardDrive struct{}

func (h *HardDrive) Read() string { return "" }

type ComputerFacade struct {
	cpu       CPU
	harddrive HardDrive
}

func (c *ComputerFacade) Start() {
	c.cpu.Freeze()
	c.harddrive.Read()
	c.cpu.Execute()
}

// func main() {

// 	cf := ComputerFacade{}
// 	cf.Start()
// }

// Применимость паттерна Фасад

//     Сложная подсистема
//     Разделение уровней системы
//     Упрощение интеграции

// Плюсы

//     Упрощение использования
//     Снижение связанности
//     Инкапсуляция подсистемы
//     Повышение читабельности и удобства поддержки

// Минусы

//     Ограничение гибкости
//     Сложность при избыточном использовании
//     Скрытие деталей реализации
