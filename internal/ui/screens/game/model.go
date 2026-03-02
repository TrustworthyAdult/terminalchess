package game

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/notnil/chess"

	"terminalchess/internal/ui/board"
	"terminalchess/internal/ui/navigate"
	"terminalchess/internal/ui/styles"
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Select key.Binding
	Flip   key.Binding
	Back   key.Binding
	Quit   key.Binding
	Help   key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Select, k.Flip, k.Back, k.Help}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Select, k.Flip, k.Back, k.Quit},
	}
}

var keys = keyMap{
	Up:     key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	Down:   key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	Left:   key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "left")),
	Right:  key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "right")),
	Select: key.NewBinding(key.WithKeys("enter", " "), key.WithHelp("enter", "select/move")),
	Flip:   key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "flip board")),
	Back:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	Quit:   key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
	Help:   key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "more")),
}

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
	help        help.Model
}

func NewModel(p Props) Model {
	return Model{
		styles:      p.Styles,
		game:        chess.NewGame(),
		cursor:      chess.A1,
		perspective: chess.White,
		help:        help.New(),
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(k, keys.Quit):
			return m, tea.Quit
		case key.Matches(k, keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(k, keys.Back):
			if m.selected != nil {
				m.selected = nil
				m.validDests = nil
			} else {
				return m, navigate.To(navigate.Menu)
			}
		case key.Matches(k, keys.Flip):
			if m.perspective == chess.White {
				m.perspective = chess.Black
			} else {
				m.perspective = chess.White
			}
		case key.Matches(k, keys.Up):
			r := m.cursor.Rank()
			if m.perspective == chess.White && r < chess.Rank8 {
				m.cursor = chess.NewSquare(m.cursor.File(), r+1)
			} else if m.perspective == chess.Black && r > chess.Rank1 {
				m.cursor = chess.NewSquare(m.cursor.File(), r-1)
			}
		case key.Matches(k, keys.Down):
			r := m.cursor.Rank()
			if m.perspective == chess.White && r > chess.Rank1 {
				m.cursor = chess.NewSquare(m.cursor.File(), r-1)
			} else if m.perspective == chess.Black && r < chess.Rank8 {
				m.cursor = chess.NewSquare(m.cursor.File(), r+1)
			}
		case key.Matches(k, keys.Left):
			f := m.cursor.File()
			if m.perspective == chess.White && f > chess.FileA {
				m.cursor = chess.NewSquare(f-1, m.cursor.Rank())
			} else if m.perspective == chess.Black && f < chess.FileH {
				m.cursor = chess.NewSquare(f+1, m.cursor.Rank())
			}
		case key.Matches(k, keys.Right):
			f := m.cursor.File()
			if m.perspective == chess.White && f < chess.FileH {
				m.cursor = chess.NewSquare(f+1, m.cursor.Rank())
			} else if m.perspective == chess.Black && f > chess.FileA {
				m.cursor = chess.NewSquare(f-1, m.cursor.Rank())
			}
		case key.Matches(k, keys.Select):
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
	isCheckmate := m.game.Outcome() != chess.NoOutcome && m.game.Method() == chess.Checkmate

	status := lipgloss.NewStyle().Height(1).Render(statusIndicator(m.game, s))
	boardView := board.Render(m.game.Position(), s.Board, opts)
	ti := ""
	if !isCheckmate {
		ti = turnIndicator(m.game, s)
	}
	turnView := lipgloss.NewStyle().Height(1).Render(ti)
	helpView := lipgloss.NewStyle().Height(2).Render(m.help.View(keys))
	return lipgloss.JoinVertical(lipgloss.Center, status, boardView, "", turnView, "", helpView)
}

func turnIndicator(g *chess.Game, s styles.Styles) string {
	base := s.Body.Background(lipgloss.Color("#6b6b6b")).Padding(0, 1)
	if g.Position().Turn() == chess.White {
		return base.Foreground(lipgloss.Color("#FFFFFF")).Render("♚  White to move")
	}
	return base.Foreground(lipgloss.Color("#1a1a1a")).Render("♚  Black to move")
}

func statusIndicator(g *chess.Game, s styles.Styles) string {
	if g.Outcome() != chess.NoOutcome && g.Method() == chess.Checkmate {
		winner := "White"
		if g.Outcome() == chess.BlackWon {
			winner = "Black"
		}
		style := s.Body.
			Background(lipgloss.Color("#3B0000")).
			Foreground(lipgloss.Color("#C8A0A0")).
			Padding(0, 1).
			Bold(true)
		return style.Render("♚  Checkmate — " + winner + " wins")
	}

	moves := g.Moves()
	if len(moves) > 0 && moves[len(moves)-1].HasTag(chess.Check) {
		side := "White"
		if g.Position().Turn() == chess.Black {
			side = "Black"
		}
		style := s.Body.
			Background(lipgloss.Color("#CC4400")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1).
			Bold(true)
		return style.Render("♚  " + side + " is in check!")
	}

	return ""
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
