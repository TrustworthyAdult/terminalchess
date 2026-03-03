package actor

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/notnil/chess"
)

// MoveMsg is delivered when a player has chosen a move.
type MoveMsg struct {
	Move *chess.Move
	Err  error
}

// ReadyMsg is delivered when a player has finished initializing.
// Player holds the fully initialized player, or nil if Err != nil.
type ReadyMsg struct {
	Color  chess.Color
	Player Player
	Err    error
}

// Player represents one side in a chess game.
type Player interface {
	// Init performs any asynchronous setup (e.g., starting an engine).
	// Sends a ReadyMsg when complete. Returns nil if no setup is needed.
	Init(color chess.Color) tea.Cmd

	// RequestMove asks the player to choose a move for the given position.
	// Returns a Cmd that sends a MoveMsg, or nil for local/interactive players.
	RequestMove(pos *chess.Position) tea.Cmd

	// IsLocal returns true if this player interacts via keyboard in this session.
	IsLocal() bool

	// Close releases any held resources.
	Close()
}
