package root

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/notnil/chess"

	"terminalchess/internal/ui/actor"
	"terminalchess/internal/ui/navigate"
	"terminalchess/internal/ui/screens/game"
	"terminalchess/internal/ui/screens/menu"
	"terminalchess/internal/ui/screens/setup"
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
	m.current = m.makeScreen(navigate.Msg{To: navigate.Menu})
	return m
}

func (m Model) Init() tea.Cmd { return m.current.Init() }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if nav, ok := msg.(navigate.Msg); ok {
		m.current = m.makeScreen(nav)
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

func (m Model) makeScreen(nav navigate.Msg) tea.Model {
	switch nav.To {
	case navigate.Setup:
		return setup.NewModel(setup.Props{Styles: m.props.Styles})
	case navigate.Game:
		white, black := actor.Player(actor.Local{}), actor.Player(actor.Local{})
		if cfg, ok := nav.Config.(game.Config); ok && cfg.ComputerColor != nil {
			computer := actor.NewComputer(cfg.SkillLevel)
			if *cfg.ComputerColor == chess.White {
				white = computer
			} else {
				black = computer
			}
		}
		return game.NewModel(game.Props{Styles: m.props.Styles, White: white, Black: black})
	default: // navigate.Menu
		return menu.NewModel(menu.Props{Styles: m.props.Styles})
	}
}
