package navigate

import tea "github.com/charmbracelet/bubbletea"

type Screen int

const (
	Menu Screen = iota
	Game
	TermInfo
)

type Msg struct{ To Screen }

func To(s Screen) tea.Cmd {
	return func() tea.Msg { return Msg{To: s} }
}
