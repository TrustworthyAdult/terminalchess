package root

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"terminalchess/internal/ui/navigate"
	"terminalchess/internal/ui/screens/game"
	"terminalchess/internal/ui/screens/menu"
	"terminalchess/internal/ui/styles"
)

type Props struct {
	Width  int
	Height int
	Styles styles.Styles
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

func (m Model) View() string {
	return lipgloss.Place(
		m.props.Width,
		m.props.Height,
		lipgloss.Center,
		lipgloss.Center,
		m.current.View(),
	)
}

func (m Model) makeScreen(s navigate.Screen) tea.Model {
	switch s {
	case navigate.Game:
		return game.NewModel(game.Props{Styles: m.props.Styles})
	default: // navigate.Menu
		return menu.NewModel(menu.Props{Styles: m.props.Styles})
	}
}
