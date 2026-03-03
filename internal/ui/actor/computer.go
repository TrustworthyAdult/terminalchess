package actor

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
)

// Computer is a Stockfish-backed engine player.
type Computer struct {
	skillLevel int
	engine     *uci.Engine
}

// NewComputer creates a Computer player at the given Stockfish skill level (0–20).
func NewComputer(skillLevel int) Computer {
	return Computer{skillLevel: skillLevel}
}

// Init starts Stockfish asynchronously and sends a ReadyMsg when ready.
func (c Computer) Init(color chess.Color) tea.Cmd {
	skillLevel := c.skillLevel
	return func() tea.Msg {
		eng, err := uci.New("stockfish")
		if err != nil {
			return ReadyMsg{Color: color, Err: err}
		}
		_ = eng.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdUCINewGame)
		_ = eng.Run(uci.CmdSetOption{Name: "Skill Level", Value: fmt.Sprintf("%d", skillLevel)})
		return ReadyMsg{Color: color, Player: Computer{skillLevel: skillLevel, engine: eng}}
	}
}

// RequestMove asks the engine to pick a move and sends a MoveMsg.
func (c Computer) RequestMove(pos *chess.Position) tea.Cmd {
	eng := c.engine
	skillLevel := c.skillLevel
	return func() tea.Msg {
		err := eng.Run(
			uci.CmdPosition{Position: pos},
			uci.CmdGo{MoveTime: skillToMoveTime(skillLevel)},
		)
		if err != nil {
			return MoveMsg{Err: err}
		}
		return MoveMsg{Move: eng.SearchResults().BestMove}
	}
}

func (Computer) IsLocal() bool { return false }

func (c Computer) Close() {
	if c.engine != nil {
		c.engine.Close()
	}
}

func skillToMoveTime(skill int) time.Duration {
	switch {
	case skill <= 7:
		return 200 * time.Millisecond
	case skill <= 14:
		return 800 * time.Millisecond
	default:
		return 2000 * time.Millisecond
	}
}
