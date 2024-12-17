package pattern

/*
	Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern

Use Cases

    When the exact class of the object to be created isn’t known beforehand.
    When you want to decouple object creation from client code.
    Example: Plugins, database drivers, or transport services where the type of object depends on runtime conditions.

Advantages

    Loose coupling: The client code doesn’t depend on concrete implementations.
    Extensibility: New products can be added without modifying the existing code.

Drawbacks

    Complexity: Adds extra classes or interfaces.
    Overkill: May be unnecessary if the object creation logic is simple.
*/

// Transport is the product interface
type Transport interface {
	Deliver() string
}

// Car is a concrete product
type Car struct{}

func (c *Car) Deliver() string {
	return "Delivery by car."
}

// Bike is a concrete product
type Bike struct{}

func (b *Bike) Deliver() string {
	return "Delivery by bike."
}

// TransportFactory is the creator interface
type TransportFactory interface {
	CreateTransport() Transport
}

// CarFactory creates Car objects
type CarFactory struct{}

func (cf *CarFactory) CreateTransport() Transport {
	return &Car{}
}

// BikeFactory creates Bike objects
type BikeFactory struct{}

func (bf *BikeFactory) CreateTransport() Transport {
	return &Bike{}
}

// func main() {
// 	// Car factory
// 	carFactory := &CarFactory{}
// 	car := carFactory.CreateTransport()
// 	fmt.Println(car.Deliver()) // Output: Delivery by car.

// 	// Bike factory
// 	bikeFactory := &BikeFactory{}
// 	bike := bikeFactory.CreateTransport()
// 	fmt.Println(bike.Deliver()) // Output: Delivery by bike.
// }
