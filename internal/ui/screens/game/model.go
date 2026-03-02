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
	styles      styles.Styles
	game        *chess.Game
	cursor      chess.Square
	selected    *chess.Square
	validDests  map[chess.Square]bool
	perspective chess.Color
}

func NewModel(p Props) Model {
	return Model{
		styles:      p.Styles,
		game:        chess.NewGame(),
		cursor:      chess.A1,
		perspective: chess.White,
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.selected != nil {
				m.selected = nil
				m.validDests = nil
			} else {
				return m, navigate.To(navigate.Menu)
			}
		case "r":
			if m.perspective == chess.White {
				m.perspective = chess.Black
			} else {
				m.perspective = chess.White
			}
		case "up", "k":
			r := m.cursor.Rank()
			if m.perspective == chess.White && r < chess.Rank8 {
				m.cursor = chess.NewSquare(m.cursor.File(), r+1)
			} else if m.perspective == chess.Black && r > chess.Rank1 {
				m.cursor = chess.NewSquare(m.cursor.File(), r-1)
			}
		case "down", "j":
			r := m.cursor.Rank()
			if m.perspective == chess.White && r > chess.Rank1 {
				m.cursor = chess.NewSquare(m.cursor.File(), r-1)
			} else if m.perspective == chess.Black && r < chess.Rank8 {
				m.cursor = chess.NewSquare(m.cursor.File(), r+1)
			}
		case "left", "h":
			f := m.cursor.File()
			if m.perspective == chess.White && f > chess.FileA {
				m.cursor = chess.NewSquare(f-1, m.cursor.Rank())
			} else if m.perspective == chess.Black && f < chess.FileH {
				m.cursor = chess.NewSquare(f+1, m.cursor.Rank())
			}
		case "right", "l":
			f := m.cursor.File()
			if m.perspective == chess.White && f < chess.FileH {
				m.cursor = chess.NewSquare(f+1, m.cursor.Rank())
			} else if m.perspective == chess.Black && f > chess.FileA {
				m.cursor = chess.NewSquare(f-1, m.cursor.Rank())
			}
		case "enter", " ":
			if m.selected == nil {
				m.selected, m.validDests = trySelect(m.game, m.cursor)
			} else if m.validDests[m.cursor] {
				executeMove(m.game, *m.selected, m.cursor)
				m.selected = nil
				m.validDests = nil
				m.perspective = m.game.Position().Turn()
			} else {
				// Try to reselect another own piece; deselects if cursor is elsewhere.
				m.selected, m.validDests = trySelect(m.game, m.cursor)
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	s := m.styles
	opts := board.RenderOptions{
		Cursor:     m.cursor,
		Selected:   m.selected,
		ValidDests: m.validDests,
		FlipBoard:  m.perspective == chess.Black,
	}
	boardView := board.Render(m.game.Position(), s.Board, opts)
	turnIndicator := turnIndicator(m.game, s)
	hint := s.Hint.Render("arrows/hjkl  move    enter/space  select·move    r  flip    esc  back")
	return lipgloss.JoinVertical(lipgloss.Center, boardView, "", turnIndicator, "", hint)
}

func turnIndicator(g *chess.Game, s styles.Styles) string {
	base := s.Body.Background(lipgloss.Color("#6b6b6b")).Padding(0, 1)
	if g.Position().Turn() == chess.White {
		return base.Foreground(lipgloss.Color("#FFFFFF")).Render("♚  White to move")
	}
	return base.Foreground(lipgloss.Color("#1a1a1a")).Render("♚  Black to move")
}

// trySelect returns a selection and valid destinations if sq holds a piece
// belonging to the side to move with at least one legal move; nil otherwise.
func trySelect(g *chess.Game, sq chess.Square) (*chess.Square, map[chess.Square]bool) {
	piece := g.Position().Board().Piece(sq)
	if piece == chess.NoPiece || piece.Color() != g.Position().Turn() {
		return nil, nil
	}
	dests := computeValidDests(g, sq)
	if len(dests) == 0 {
		return nil, nil
	}
	return &sq, dests
}

// executeMove plays the move from→to, promoting to a queen when required.
func executeMove(g *chess.Game, from, to chess.Square) {
	for _, move := range g.ValidMoves() {
		if move.S1() == from && move.S2() == to {
			if move.Promo() == chess.NoPieceType || move.Promo() == chess.Queen {
				g.Move(move)
				return
			}
		}
	}
}

func computeValidDests(g *chess.Game, sq chess.Square) map[chess.Square]bool {
	dests := make(map[chess.Square]bool)
	for _, move := range g.ValidMoves() {
		if move.S1() == sq {
			dests[move.S2()] = true
		}
	}
	return dests
}
