package actor

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/notnil/chess"
)

// Local is a human player using keyboard input in this SSH session.
type Local struct{}

func (Local) Init(chess.Color) tea.Cmd          { return nil }
func (Local) RequestMove(*chess.Position) tea.Cmd { return nil }
func (Local) IsLocal() bool                       { return true }
func (Local) Close()                              {}
