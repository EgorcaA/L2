package pattern

import "fmt"

/*
	Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern

Benefits of Command Pattern

    Decouples sender and receiver: The invoker (e.g., RemoteControl) doesn't need to know about the receiver's implementation.
    Extensibility: Easily add new commands without changing existing code.

Drawbacks of Command Pattern

    Increased complexity: Requires additional classes for each command, which may be overkill for simple scenarios.
    Scalability concerns: Adding many commands for each receiver can clutter the codebase.

Use Cases of Command Pattern

    Undo/Redo functionality
    Task queues: Scheduling and executing tasks, like job queues.
    Remote controls: Triggering actions without tightly coupling to the receiver (e.g., home automation).
	This approach lets introduce new commands into the app without breaking any existing code.
*/

type Command interface {
	Execute()
}

// Receiver: Light
type Light struct{}

func (l *Light) On() {
	fmt.Println("The light is ON")
}

func (l *Light) Off() {
	fmt.Println("The light is OFF")
}

// Concrete Command: TurnOnLight
type TurnOnLightCommand struct {
	light *Light
}

func (c *TurnOnLightCommand) Execute() {
	c.light.On()
}

// Concrete Command: TurnOffLight
type TurnOffLightCommand struct {
	light *Light
}

func (c *TurnOffLightCommand) Execute() {
	c.light.Off()
}

// Invoker
type RemoteControl struct {
	command Command
}

func (r *RemoteControl) SetCommand(command Command) {
	r.command = command
}

func (r *RemoteControl) PressButton() {
	r.command.Execute()
}

// func main() {
// 	// Receiver: Light
// 	light := &Light{}

// 	// Concrete commands
// 	turnOnLight := &TurnOnLightCommand{light: light}
// 	turnOffLight := &TurnOffLightCommand{light: light}

// 	// Invoker: RemoteControl
// 	remote := &RemoteControl{}

// 	// Turn the light ON
// 	remote.SetCommand(turnOnLight)
// 	remote.PressButton()

// 	// Turn the light OFF
// 	remote.SetCommand(turnOffLight)
// 	remote.PressButton()
// }
