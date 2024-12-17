package pattern

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern

Use Cases

    Event handling: Chains of event listeners.
    Logging systems: Handling logs at different levels (e.g., Debug, Info, Error).
    Request filtering: Processing requests (e.g., middleware in web servers).
    Approval workflows: Sequential approval processes in organizations.

Benefits
	You can control the order of request handling.
	Single Responsibility Principle. You can decouple classes that invoke operations from classes that perform operations.
	Open/Closed Principle. You can introduce new handlers into the app without breaking the existing client code.

Drawbacks

    Debugging complexity: Hard to track the flow if the chain is long.
    No guarantee of handling: If no handler processes the request, it might be dropped silently.
    Performance overhead: Requests may traverse many handlers unnecessarily.
*/

import "fmt"

// Handler interface
type Handler interface {
	SetNext(handler Handler) Handler
	Handle(request string)
}

// BaseHandler provides default behavior for chaining
type BaseHandler struct {
	next Handler
}

func (b *BaseHandler) SetNext(handler Handler) Handler {
	b.next = handler
	return handler
}

func (b *BaseHandler) Handle(request string) {
	if b.next != nil {
		b.next.Handle(request)
	}
}

// DebugHandler
type DebugHandler struct {
	BaseHandler
}

func (d *DebugHandler) Handle(request string) {
	if request == "debug" {
		fmt.Println("DebugHandler: Handling debug log")
	} else {
		fmt.Println("DebugHandler: Passing to next")
		d.BaseHandler.Handle(request)
	}
}

// InfoHandler
type InfoHandler struct {
	BaseHandler
}

func (i *InfoHandler) Handle(request string) {
	if request == "info" {
		fmt.Println("InfoHandler: Handling info log")
	} else {
		fmt.Println("InfoHandler: Passing to next")
		i.BaseHandler.Handle(request) // WOW
	}
}

// ErrorHandler
type ErrorHandler struct {
	BaseHandler
}

func (e *ErrorHandler) Handle(request string) {
	if request == "error" {
		fmt.Println("ErrorHandler: Handling error log")
	} else {
		fmt.Println("ErrorHandler: Passing to next")
		e.BaseHandler.Handle(request)
	}
}

// func main() {
// 	// Create handlers
// 	debugHandler := &DebugHandler{}
// 	infoHandler := &InfoHandler{}
// 	errorHandler := &ErrorHandler{}

// 	// Build the chain: Debug -> Info -> Error
// 	debugHandler.SetNext(infoHandler).SetNext(errorHandler)

// 	// Send requests
// 	fmt.Println("Sending 'debug' request:")
// 	debugHandler.Handle("debug")

// 	fmt.Println("\nSending 'info' request:")
// 	debugHandler.Handle("info")

// 	fmt.Println("\nSending 'error' request:")
// 	debugHandler.Handle("error")

// 	fmt.Println("\nSending 'unknown' request:")
// 	debugHandler.Handle("unknown")
// }
