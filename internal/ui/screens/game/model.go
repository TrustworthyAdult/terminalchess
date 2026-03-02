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

type panelFocus int

const (
	boardFocus panelFocus = iota
	moveListFocus
	chatFocus
)

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Select   key.Binding
	Flip     key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
	Back     key.Binding
	Quit     key.Binding
	Help     key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Select, k.Tab, k.Back, k.Help}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Select, k.Flip, k.Tab, k.ShiftTab, k.Back, k.Quit},
	}
}

var keys = keyMap{
	Up:       key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	Down:     key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	Left:     key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "left")),
	Right:    key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "right")),
	Select:   key.NewBinding(key.WithKeys("enter", " "), key.WithHelp("enter", "select/move")),
	Flip:     key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "flip board")),
	Tab:      key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next panel")),
	ShiftTab: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev panel")),
	Back:     key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	Quit:     key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
	Help:     key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "more")),
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
	focus       panelFocus
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
		case key.Matches(k, keys.Tab):
			m.focus = (m.focus + 1) % 3
		case key.Matches(k, keys.ShiftTab):
			m.focus = (m.focus + 2) % 3
		case key.Matches(k, keys.Back):
			if m.selected != nil {
				m.selected = nil
				m.validDests = nil
			} else {
				return m, navigate.To(navigate.Menu)
			}
		case m.focus == boardFocus && key.Matches(k, keys.Flip):
			if m.perspective == chess.White {
				m.perspective = chess.Black
			} else {
				m.perspective = chess.White
			}
		case m.focus == boardFocus && key.Matches(k, keys.Up):
			r := m.cursor.Rank()
			if m.perspective == chess.White && r < chess.Rank8 {
				m.cursor = chess.NewSquare(m.cursor.File(), r+1)
			} else if m.perspective == chess.Black && r > chess.Rank1 {
				m.cursor = chess.NewSquare(m.cursor.File(), r-1)
			}
		case m.focus == boardFocus && key.Matches(k, keys.Down):
			r := m.cursor.Rank()
			if m.perspective == chess.White && r > chess.Rank1 {
				m.cursor = chess.NewSquare(m.cursor.File(), r-1)
			} else if m.perspective == chess.Black && r < chess.Rank8 {
				m.cursor = chess.NewSquare(m.cursor.File(), r+1)
			}
		case m.focus == boardFocus && key.Matches(k, keys.Left):
			f := m.cursor.File()
			if m.perspective == chess.White && f > chess.FileA {
				m.cursor = chess.NewSquare(f-1, m.cursor.Rank())
			} else if m.perspective == chess.Black && f < chess.FileH {
				m.cursor = chess.NewSquare(f+1, m.cursor.Rank())
			}
		case m.focus == boardFocus && key.Matches(k, keys.Right):
			f := m.cursor.File()
			if m.perspective == chess.White && f < chess.FileH {
				m.cursor = chess.NewSquare(f+1, m.cursor.Rank())
			} else if m.perspective == chess.Black && f > chess.FileA {
				m.cursor = chess.NewSquare(f-1, m.cursor.Rank())
			}
		case m.focus == boardFocus && key.Matches(k, keys.Select):
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

func panelBorderStyle(focused bool) lipgloss.Style {
	color := lipgloss.Color("#444444")
	if focused {
		color = lipgloss.Color("#FFFFFF")
	}
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(color).
		Padding(1, 2)
}

func (m Model) View() string {
	s := m.styles
	opts := board.RenderOptions{
		Cursor:     m.cursor,
		Selected:   m.selected,
		ValidDests: m.validDests,
		FlipBoard:  m.perspective == chess.Black,
	}
	indicator := lipgloss.NewStyle().Height(1).Render(gameIndicator(m.game, s))
	boardView := board.Render(m.game.Position(), s.Board, opts)

	boardContent := lipgloss.JoinVertical(lipgloss.Center, boardView, indicator)
	boardPanel := panelBorderStyle(m.focus == boardFocus).Render(boardContent)

	boardH := lipgloss.Height(boardPanel)
	const vOverhead = 4 // 2 border rows + 2 padding rows per panel
	innerH := boardH - vOverhead
	if innerH < 0 {
		innerH = 0
	}

	chatContent := lipgloss.NewStyle().Height(innerH).Render("Chat")
	moveHistoryContent := lipgloss.NewStyle().Height(innerH).Render("Move History")
	chatPanel := panelBorderStyle(m.focus == chatFocus).Render(chatContent)
	moveHistoryPanel := panelBorderStyle(m.focus == moveListFocus).Render(moveHistoryContent)

	panels := lipgloss.JoinHorizontal(lipgloss.Top, chatPanel, boardPanel, moveHistoryPanel)
	helpView := lipgloss.NewStyle().Height(2).Render(m.help.View(keys))
	return lipgloss.JoinVertical(lipgloss.Left, panels, helpView)
}

func gameIndicator(g *chess.Game, s styles.Styles) string {
	base := s.Body.Padding(0, 1)

	if g.Outcome() != chess.NoOutcome && g.Method() == chess.Checkmate {
		winner := "White"
		if g.Outcome() == chess.BlackWon {
			winner = "Black"
		}
		return base.
			Background(lipgloss.Color("#3B0000")).
			Foreground(lipgloss.Color("#C8A0A0")).
			Bold(true).
			Render("♚  Checkmate — " + winner + " wins")
	}

	moves := g.Moves()
	if len(moves) > 0 && moves[len(moves)-1].HasTag(chess.Check) {
		side := "White"
		if g.Position().Turn() == chess.Black {
			side = "Black"
		}
		return base.
			Background(lipgloss.Color("#CC4400")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Render("♚  " + side + " is in check!")
	}

	if g.Position().Turn() == chess.White {
		return base.
			Background(lipgloss.Color("#6b6b6b")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Render("♚  White to move")
	}
	return base.
		Background(lipgloss.Color("#6b6b6b")).
		Foreground(lipgloss.Color("#1a1a1a")).
		Render("♚  Black to move")
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
