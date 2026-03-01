package game

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/notnil/chess"

	"terminalchess/internal/ui/board"
	"terminalchess/internal/ui/navigate"
	"terminalchess/internal/ui/styles"
)

type Props struct {
	Styles styles.Styles
}

type Model struct {
	styles styles.Styles
	game   *chess.Game
}

func NewModel(p Props) Model {
	return Model{
		styles: p.Styles,
		game:   chess.NewGame(),
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "q", "esc":
			return m, navigate.To(navigate.Menu)
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	s := m.styles
	boardView := board.Render(m.game.Position(), s.Board)
	hint := s.Hint.Render("esc  go back")
	return lipgloss.JoinVertical(lipgloss.Center, boardView, "", hint)
}
