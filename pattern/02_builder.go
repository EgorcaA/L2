package pattern

import "fmt"

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern
*/

type Bicycle struct {
	Color string
	Model string
}

// The builder abstraction.
type IBicycleBuilder interface {
	SetColor(color string)
	MakeBicycle() *Bicycle
}

// Concrete builder: GTBuilder
type GTBuilder struct {
	color string
}

func (bld *GTBuilder) SetColor(color string) {
	bld.color = color
}

func (bld *GTBuilder) MakeBicycle() *Bicycle {
	return &Bicycle{Color: bld.color, Model: "GT"}
}

// The director.
type MountainBikeBuildDirector struct {
	color   string
	builder IBicycleBuilder
}

func (director *MountainBikeBuildDirector) Construct() {
	director.builder.SetColor(director.color)
}

func (director *MountainBikeBuildDirector) GetBicycle() *Bicycle {
	return director.builder.MakeBicycle()
}

func main() {
	// Director controls the stepwise creation of product and returns the result.
	builder := &GTBuilder{}
	director := MountainBikeBuildDirector{
		color:   "red",
		builder: builder,
	}

	director.Construct()
	bc := director.GetBicycle()
	fmt.Println(bc.Color)
	fmt.Println(bc.Model)
}

// Pluses of Builder Pattern

//     Step-by-step construction: Allows the creation of complex objects in a controlled, incremental manner.
//     Code readability: Separates object construction logic from its representation, improving clarity.
//     Flexibility: Easily create variations of an object by reusing or modifying steps.
//     Encapsulation: Hides the construction process from the client, exposing only the finished product.
//     Immutable objects: Supports creation of immutable objects by ensuring all required fields are initialized.

// Minuses of Builder Pattern

//     Overhead for simple objects: Adds unnecessary complexity for straightforward object creation.
//     Increased codebase size: Requires additional classes for builders and directors.
//     Not always intuitive: May be harder to understand compared to simpler patterns.

// Use Cases for Builder Pattern

//     Complex object creation: When an object requires multiple configuration steps (e.g., constructing UI components or assembling a vehicle).
//     Multiple representations: When the same construction process should produce different types of objects (e.g., different document formats).
//     Immutable objects: When objects should be immutable after construction.
//     Reducing constructor overload: When there are too many parameters in a constructor, the builder pattern can improve readability and maintainability.
