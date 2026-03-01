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
	cursor chess.Square
}

func NewModel(p Props) Model {
	return Model{
		styles: p.Styles,
		game:   chess.NewGame(),
		cursor: chess.A1,
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
		case "up", "k":
			if r := m.cursor.Rank(); r < chess.Rank8 {
				m.cursor = chess.NewSquare(m.cursor.File(), r+1)
			}
		case "down", "j":
			if r := m.cursor.Rank(); r > chess.Rank1 {
				m.cursor = chess.NewSquare(m.cursor.File(), r-1)
			}
		case "left", "h":
			if f := m.cursor.File(); f > chess.FileA {
				m.cursor = chess.NewSquare(f-1, m.cursor.Rank())
			}
		case "right", "l":
			if f := m.cursor.File(); f < chess.FileH {
				m.cursor = chess.NewSquare(f+1, m.cursor.Rank())
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	s := m.styles
	boardView := board.Render(m.game.Position(), s.Board, m.cursor)
	hint := s.Hint.Render("arrows/hjkl  move    esc  go back")
	return lipgloss.JoinVertical(lipgloss.Center, boardView, "", hint)
}
