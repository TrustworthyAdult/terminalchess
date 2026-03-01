package root

import (
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"

	"terminalchess/internal/ui/navigate"
	"terminalchess/internal/ui/screens/game"
	"terminalchess/internal/ui/screens/menu"
)

type Props struct {
	Width     int
	Height    int
	TxtStyle  lipgloss.Style
	QuitStyle lipgloss.Style
}

type Model struct {
	props   Props
	current tea.Model
}

func New(p Props) Model {
	m := Model{props: p}
	m.current = m.makeScreen(navigate.Menu)
	return m
}

func (m Model) Init() tea.Cmd { return m.current.Init() }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if nav, ok := msg.(navigate.Msg); ok {
		m.current = m.makeScreen(nav.To)
		return m, m.current.Init()
	}

	if wsm, ok := msg.(tea.WindowSizeMsg); ok {
		m.props.Width = wsm.Width
		m.props.Height = wsm.Height
	}

	var cmd tea.Cmd
	m.current, cmd = m.current.Update(msg)
	return m, cmd
}

func (m Model) View() string { return m.current.View() }

func (m Model) makeScreen(s navigate.Screen) tea.Model {
	p := m.props
	switch s {
	case navigate.Game:
		return game.NewModel(game.Props{
			TxtStyle:  p.TxtStyle,
			QuitStyle: p.QuitStyle,
		})
	default: // navigate.Menu
		return menu.NewModel(menu.Props{
			Width:     p.Width,
			Height:    p.Height,
			TxtStyle:  p.TxtStyle,
			QuitStyle: p.QuitStyle,
		})
	}
}
