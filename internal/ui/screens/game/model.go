package game

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	ansi "github.com/charmbracelet/x/ansi"
	"github.com/notnil/chess"

	"terminalchess/internal/ui/board"
	"terminalchess/internal/ui/navigate"
	"terminalchess/internal/ui/styles"
)

type gamePhase int

const (
	playing  gamePhase = iota
	postGame           // popup visible
	studying           // popup dismissed, board read-only
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
	moveHistory viewport.Model
	termWidth   int
	phase       gamePhase
	popupChoice int
}

func NewModel(p Props) Model {
	return Model{
		styles:      p.Styles,
		game:        chess.NewGame(),
		cursor:      chess.A1,
		perspective: chess.White,
		help:        help.New(),
		moveHistory: newMoveHistoryViewport(p.Styles),
	}
}

// newMoveHistoryViewport creates a viewport sized to match the board panel
// height. The height is computed from a dummy render so GotoBottom works
// correctly in Update before the first View call.
func newMoveHistoryViewport(s styles.Styles) viewport.Model {
	dummyBoard := board.Render(chess.NewGame().Position(), s.Board, board.RenderOptions{})
	dummyInd := lipgloss.NewStyle().Height(1).Render("")
	dummyContent := lipgloss.JoinVertical(lipgloss.Center, dummyBoard, dummyInd)
	dummyPanel := panelBorderStyle(false).Render(dummyContent)
	const vOverhead = 4
	innerH := lipgloss.Height(dummyPanel) - vOverhead
	if innerH < 2 {
		innerH = 10
	}
	return viewport.New(20, innerH-1) // -1 for title row
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if wsm, ok := msg.(tea.WindowSizeMsg); ok {
		m.termWidth = wsm.Width
	}

	if k, ok := msg.(tea.KeyMsg); ok {
		// Quit always works regardless of phase.
		if key.Matches(k, keys.Quit) {
			return m, tea.Quit
		}

		if m.phase == postGame {
			switch {
			case key.Matches(k, keys.Up):
				m.popupChoice = (m.popupChoice + 2) % 3
			case key.Matches(k, keys.Down):
				m.popupChoice = (m.popupChoice + 1) % 3
			case key.Matches(k, keys.Select):
				switch m.popupChoice {
				case 0: // Exit
					return m, navigate.To(navigate.Menu)
				case 1: // Rematch
					m.game = chess.NewGame()
					m.cursor = chess.A1
					m.selected = nil
					m.validDests = nil
					m.perspective = chess.White
					m.moveHistory = newMoveHistoryViewport(m.styles)
					m.phase = playing
					m.popupChoice = 0
				case 2: // Study
					m.phase = studying
				}
			case key.Matches(k, keys.Back):
				m.phase = studying
			}
			// All other keys consumed in postGame.
			return m, nil
		}

		// playing or studying
		switch {
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
		case m.focus == moveListFocus && key.Matches(k, keys.Up):
			m.moveHistory.ScrollUp(1)
		case m.focus == moveListFocus && key.Matches(k, keys.Down):
			m.moveHistory.ScrollDown(1)
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
		case m.focus == boardFocus && m.phase == playing && key.Matches(k, keys.Select):
			if m.selected == nil {
				m.selected, m.validDests = trySelect(m.game, m.cursor)
			} else if m.validDests[m.cursor] {
				executeMove(m.game, *m.selected, m.cursor)
				m.selected = nil
				m.validDests = nil
				m.perspective = m.game.Position().Turn()
				m.moveHistory.SetContent(strings.Join(moveHistoryLines(m.game), "\n"))
				m.moveHistory.GotoBottom()
				if m.game.Outcome() != chess.NoOutcome {
					m.phase = postGame
				}
			} else {
				// Try to reselect another own piece; deselects if cursor is elsewhere.
				m.selected, m.validDests = trySelect(m.game, m.cursor)
			}
		}
	}
	return m, nil
}

func outcomeText(g *chess.Game) string {
	outcome := g.Outcome()
	method := g.Method()
	if method == chess.Checkmate {
		if outcome == chess.WhiteWon {
			return "Checkmate — White wins"
		}
		return "Checkmate — Black wins"
	}
	switch method {
	case chess.Stalemate:
		return "Draw — Stalemate"
	case chess.ThreefoldRepetition:
		return "Draw — Threefold Repetition"
	case chess.FiftyMoveRule:
		return "Draw — 50-Move Rule"
	case chess.InsufficientMaterial:
		return "Draw — Insufficient Material"
	default:
		return "Draw"
	}
}

func renderPopup(g *chess.Game, choice int) string {
	labels := []string{"Exit", "Rematch", "Study"}
	titleW := lipgloss.Width(outcomeText(g))
	selectedStyle := lipgloss.NewStyle().Bold(true).
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#ffffff")).
		Width(titleW)
	unselectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#777777")).
		Background(lipgloss.Color("#161616")).
		Width(titleW)

	var sb strings.Builder
	for i, label := range labels {
		if i == choice {
			sb.WriteString(selectedStyle.Render("▸ " + label))
		} else {
			sb.WriteString(unselectedStyle.Render("  " + label))
		}
		if i < len(labels)-1 {
			sb.WriteString("\n")
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffffff")).Render(outcomeText(g)),
		"",
		sb.String(),
	)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#ffffff")).
		Padding(1, 3).
		Background(lipgloss.Color("#161616")).
		Render(content)
}

func overlayCenter(bg, fg string) string {
	bgLines := strings.Split(bg, "\n")
	fgLines := strings.Split(fg, "\n")

	bgW := lipgloss.Width(bg)
	fgW := lipgloss.Width(fg)
	startY := (len(bgLines) - len(fgLines)) / 2
	startX := (bgW - fgW) / 2
	if startX < 0 {
		startX = 0
	}
	if startY < 0 {
		startY = 0
	}

	for i, fgLine := range fgLines {
		idx := startY + i
		if idx >= len(bgLines) {
			continue
		}
		bgLine := bgLines[idx]
		lineW := lipgloss.Width(bgLine)
		fgLineW := lipgloss.Width(fgLine)

		left := ansi.Truncate(bgLine, startX, "")
		if w := lipgloss.Width(left); w < startX {
			left += strings.Repeat(" ", startX-w)
		}
		right := ansi.Cut(bgLine, startX+fgLineW, lineW)
		bgLines[idx] = left + fgLine + right
	}
	return strings.Join(bgLines, "\n")
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

	boardPanelW := lipgloss.Width(boardPanel)
	boardH := lipgloss.Height(boardPanel)
	const vOverhead = 4 // 2 border rows + 2 padding rows per panel
	innerH := boardH - vOverhead
	if innerH < 2 {
		innerH = 2
	}

	sideInnerW := (m.termWidth - boardPanelW) / 2 - 6 // 6 = border(2) + padding(4)
	if sideInnerW < 15 {
		sideInnerW = 15
	}

	// Apply alternating row backgrounds to move history content.
	evenRow := lipgloss.NewStyle().Width(sideInnerW).Background(lipgloss.Color("#1e1e1e"))
	oddRow := lipgloss.NewStyle().Width(sideInnerW).Background(lipgloss.Color("#2a2a2a"))
	lines := moveHistoryLines(m.game)
	for i, line := range lines {
		if i%2 == 0 {
			lines[i] = evenRow.Render(line)
		} else {
			lines[i] = oddRow.Render(line)
		}
	}
	m.moveHistory.Width = sideInnerW
	m.moveHistory.Height = innerH - 1 // -1 for title row
	m.moveHistory.SetContent(strings.Join(lines, "\n"))

	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#aaaaaa")).Render("Move History")
	moveHistoryContent := lipgloss.JoinVertical(lipgloss.Left, title, m.moveHistory.View())
	moveHistoryPanel := panelBorderStyle(m.focus == moveListFocus).Render(moveHistoryContent)

	chatContent := lipgloss.NewStyle().Width(sideInnerW).Height(innerH).Render("Chat")
	chatPanel := panelBorderStyle(m.focus == chatFocus).Render(chatContent)

	panels := lipgloss.JoinHorizontal(lipgloss.Top, chatPanel, boardPanel, moveHistoryPanel)
	helpView := lipgloss.NewStyle().Height(2).Render(m.help.View(keys))
	screen := lipgloss.JoinVertical(lipgloss.Left, panels, helpView)
	if m.phase == postGame {
		return overlayCenter(screen, renderPopup(m.game, m.popupChoice))
	}
	return screen
}

func moveHistoryLines(g *chess.Game) []string {
	history := g.MoveHistory()
	an := chess.AlgebraicNotation{}
	lines := make([]string, 0, len(history)/2+1)
	for i := 0; i < len(history); i += 2 {
		white := an.Encode(history[i].PrePosition, history[i].Move)
		line := fmt.Sprintf("%d. %s", i/2+1, white)
		if i+1 < len(history) {
			line += "  " + an.Encode(history[i+1].PrePosition, history[i+1].Move)
		}
		lines = append(lines, line)
	}
	return lines
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
