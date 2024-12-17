package pattern

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern
// https://refactoring.guru/ru/design-patterns/visitor
*/

type Shape interface {
	move(x, y float32)
	draw()
	accept(v Visitor)
}

// Метод принятия посетителя должен быть реализован в каждом
// элементе, а не только в базовом классе. Это поможет программе
// определить, какой метод посетителя нужно вызвать, если вы не
// знаете тип элемента.
type Dot struct{}

func (d *Dot) move(x, y float32) {}
func (d *Dot) draw()             {}
func (d *Dot) accept(v Visitor)  { v.visitDot(d) }

type Circle struct{}

func (c *Circle) move(x, y float32) {}
func (c *Circle) draw()             {}
func (c *Circle) accept(v Visitor)  { v.visitCircle(c) }

type CompoundShape struct{}

func (cs *CompoundShape) move(x, y float32) {}
func (cs *CompoundShape) draw()             {}
func (cs *CompoundShape) accept(v Visitor)  { v.visitCompoundShape(cs) }

// Интерфейс посетителей должен содержать методы посещения
// каждого элемента. Важно, чтобы иерархия элементов менялась
// редко, так как при добавлении нового элемента придётся менять
// всех существующих посетителей.
type Visitor interface {
	visitDot(d *Dot)
	visitCircle(c *Circle)
	visitCompoundShape(cs *CompoundShape)
}

// Конкретный посетитель реализует одну операцию для всей
// иерархии элементов. Новая операция = новый посетитель.
// Посетитель выгодно применять, когда новые элементы
// добавляются очень редко, а новые операции — часто.
type XMLExportVisitor struct{}

func (vst *XMLExportVisitor) visitDot(d *Dot) {
	// Экспорт id и координат центра точки.
}
func (vst *XMLExportVisitor) visitCircle(c *Circle) {
	// Экспорт id, кординат центра и радиуса окружности.
}
func (vst *XMLExportVisitor) visitCompoundShape(cs *CompoundShape) {
	// Экспорт id составной фигуры, а также списка id
	// подфигур, из которых она состоит.
}

// Приложение может применять посетителя к любому набору
// объектов элементов, даже не уточняя их типы. Нужный метод
// посетителя будет выбран благодаря проходу через метод accept.

// func main() {
// 	allShapes := []Shape{&Dot{}, &Circle{}, &CompoundShape{}}
// 	exporter := &XMLExportVisitor{}
// 	for _, shape := range allShapes {
// 		shape.accept(exporter)
// 	}
// }
