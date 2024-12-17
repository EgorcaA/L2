package pattern

/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern

Use Cases

    When an object’s behavior depends on its state, and it needs to change its behavior dynamically.
    Finite State Machines (FSM), such as traffic lights, vending machines, or document workflows.
    Simplifying complex conditionals based on state (if or switch statements).

Advantages

    Clean State Transitions: Encapsulates state-specific behavior in separate classes.
    Open/Closed Principle: New states can be added without modifying the existing code.
    Readability: Eliminates large conditionals related to state changes.

Drawbacks

    Complexity: Adds additional classes for each state.
    Overkill: May not be needed if the state logic is simple.
*/

import "fmt"

// Player is the context that changes behavior based on its state
type Player struct {
	state State
}

// NewPlayer initializes the player with a default state
func NewPlayer() *Player {
	return &Player{state: &LockedState{}} // Starts in LockedState
}

// SetState changes the current state of the player
func (p *Player) SetState(state State) {
	p.state = state
}

// PressPlay delegates the action to the current state
func (p *Player) PressPlay() {
	p.state.PressPlay(p)
}

// PressLock delegates the action to the current state
func (p *Player) PressLock() {
	p.state.PressLock(p)
}

// State defines the interface for the player states
type State interface {
	PressPlay(player *Player)
	PressLock(player *Player)
}

// LockedState represents the state where the player is locked
type LockedState struct{}

func (l *LockedState) PressPlay(player *Player) {
	fmt.Println("Player is locked. Unlock it first.")
}

func (l *LockedState) PressLock(player *Player) {
	fmt.Println("Unlocking the player...")
	player.SetState(&PlayingState{}) // Transition to PlayingState
}

// PlayingState represents the state where the player is playing music
type PlayingState struct{}

func (p *PlayingState) PressPlay(player *Player) {
	fmt.Println("Pausing music...")
	player.SetState(&LockedState{}) // Transition to LockedState
}

func (p *PlayingState) PressLock(player *Player) {
	fmt.Println("Locking the player...")
	player.SetState(&LockedState{}) // Transition to LockedState
}

// func main() {
// 	player := NewPlayer()

// 	// Player starts in LockedState
// 	player.PressPlay() // Output: Player is locked. Unlock it first.
// 	player.PressLock() // Output: Unlocking the player...

// 	// Now in PlayingState
// 	player.PressPlay() // Output: Pausing music...
// 	player.PressLock() // Output: Locking the player...

// 	// Back to LockedState
// 	player.PressPlay() // Output: Player is locked. Unlock it first.
// }
