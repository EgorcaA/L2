package pattern

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern

Use Cases

    When you have multiple algorithms (e.g., sorting, payment methods, compression) and want to choose one at runtime.
    When you want to avoid conditional statements (e.g., if or switch) to select an algorithm.
    Example: Payment processing systems, sorting algorithms, or routing policies in network systems.

Advantages

    Flexibility: Easily switch between algorithms at runtime.
    Extensibility: Add new strategies without modifying the existing code.
    Cleaner Code: Eliminates large conditional statements.

Drawbacks

    Overhead: Adds complexity with additional classes or structs.
    Client Awareness: The client must know about different strategies to set them
*/

// The strategy interface declares operations common to all
// supported versions of some algorithm. The context uses this
// interface to call the algorithm defined by the concrete
// strategies.
type Strategy interface {
	Execute(a, b int) int
}

// Concrete strategies implement the algorithm while following
// the base strategy interface. The interface makes them
// interchangeable in the context.
type ConcreteStrategyAdd struct{}

func (add *ConcreteStrategyAdd) Execute(a, b int) int {
	return a + b
}

type ConcreteStrategySubtract struct{}

func (add *ConcreteStrategySubtract) Execute(a, b int) int {
	return a - b
}

// The context defines the interface of interest to clients.
type Context struct {
	// The context maintains a reference to one of the strategy
	// objects. The context doesn't know the concrete class of a
	// strategy. It should work with all strategies via the
	// strategy interface.
	strategy Strategy
}

// Usually the context accepts a strategy through the
// constructor, and also provides a setter so that the
// strategy can be switched at runtime.
func (ctx *Context) SetStrategy(strategy Strategy) {
	ctx.strategy = strategy
}

// The context delegates some work to the strategy object
// instead of implementing multiple versions of the
// algorithm on its own.
func (ctx *Context) ExecuteStrategy(a, b int) int {
	return ctx.strategy.Execute(a, b)
}

// The client code picks a concrete strategy and passes it to
// the context. The client should be aware of the differences
// between strategies in order to make the right choice.

// func main() {
// 	ctx := &Context{}
// 	ctx.SetStrategy(&ConcreteStrategyAdd{})
// 	result := ctx.ExecuteStrategy(3, 5)
// 	fmt.Println(result)
// }
