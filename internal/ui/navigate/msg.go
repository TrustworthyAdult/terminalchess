package navigate

import tea "github.com/charmbracelet/bubbletea"

type Screen int

const (
	Menu  Screen = iota
	Game
	Setup
)

type Msg struct {
	To     Screen
	Config any
}

func To(s Screen) tea.Cmd {
	return func() tea.Msg { return Msg{To: s} }
}

func ToWithConfig(s Screen, cfg any) tea.Cmd {
	return func() tea.Msg { return Msg{To: s, Config: cfg} }
}
