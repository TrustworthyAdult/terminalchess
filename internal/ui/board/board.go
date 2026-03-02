package board

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/notnil/chess"
)

type Styles struct {
	LightSquare    lipgloss.Style
	DarkSquare     lipgloss.Style
	CursorLight    lipgloss.Style
	CursorDark     lipgloss.Style
	SelectedLight  lipgloss.Style
	SelectedDark   lipgloss.Style
	LastMoveLight  lipgloss.Style
	LastMoveDark   lipgloss.Style
	Label          lipgloss.Style
	LabelHighlight lipgloss.Style
}

func NewStyles(renderer *lipgloss.Renderer) Styles {
	cursorLight := lipgloss.Color("#CCD277")
	cursorDark := lipgloss.Color("#AAA23A")
	return Styles{
		LightSquare:   renderer.NewStyle().Background(lipgloss.Color("#F0D9B5")),
		DarkSquare:    renderer.NewStyle().Background(lipgloss.Color("#B58863")),
		CursorLight:   renderer.NewStyle().Background(cursorLight),
		CursorDark:    renderer.NewStyle().Background(cursorDark),
		SelectedLight: renderer.NewStyle().Background(lipgloss.Color("#C8A800")),
		SelectedDark:  renderer.NewStyle().Background(lipgloss.Color("#D4AC00")),
		LastMoveLight: renderer.NewStyle().Background(lipgloss.Color("#CDD26A")),
		LastMoveDark:  renderer.NewStyle().Background(lipgloss.Color("#AABA72")),
		Label: renderer.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#9b9b9b", Dark: "#5c5c5c"}),
		LabelHighlight: renderer.NewStyle().
			Foreground(cursorLight).
			Bold(true),
	}
}

var pieces = map[chess.PieceType]string{
	chess.King:   "♚",
	chess.Queen:  "♛",
	chess.Rook:   "♜",
	chess.Bishop: "♝",
	chess.Knight: "♞",
	chess.Pawn:   "♙",
}

type RenderOptions struct {
	Cursor        chess.Square
	Selected      *chess.Square
	ValidDests    map[chess.Square]bool
	FlipBoard     bool
	LastMoveFrom  *chess.Square
	LastMoveTo    *chess.Square
}

// Render draws the board. When FlipBoard is false the view is from White's
// side (rank 1 at the bottom); when true it is from Black's side.
func Render(pos *chess.Position, s Styles, opts RenderOptions) string {
	b := pos.Board()
	cursorFile := opts.Cursor.File()
	cursorRank := opts.Cursor.Rank()
	rankOrder, fileOrder := boardOrder(opts.FlipBoard)

	var sb strings.Builder
	sb.WriteString(fileLabels(s, cursorFile, fileOrder))
	sb.WriteString("\n")

	for _, r := range rankOrder {
		rankLabel := s.labelStyle(r == cursorRank)
		sb.WriteString(rankLabel.Render(fmt.Sprintf("%d ", int(r)+1)))
		for _, f := range fileOrder {
			sq := chess.NewSquare(f, r)
			sb.WriteString(renderSquare(sq, b.Piece(sq), opts, s))
		}
		sb.WriteString(rankLabel.Render(fmt.Sprintf(" %d", int(r)+1)))
		sb.WriteString("\n")
	}

	sb.WriteString(fileLabels(s, cursorFile, fileOrder))
	return sb.String()
}

// boardOrder returns the rank and file iteration order for the given perspective.
// White: ranks 8→1 top-to-bottom, files A→H left-to-right.
// Black: ranks 1→8 top-to-bottom, files H→A left-to-right.
func boardOrder(flip bool) ([]chess.Rank, []chess.File) {
	ranks := make([]chess.Rank, 8)
	files := make([]chess.File, 8)
	for i := 0; i < 8; i++ {
		if flip {
			ranks[i] = chess.Rank(i)
			files[i] = chess.File(7 - i)
		} else {
			ranks[i] = chess.Rank(7 - i)
			files[i] = chess.File(i)
		}
	}
	return ranks, files
}

func (s Styles) labelStyle(highlighted bool) lipgloss.Style {
	if highlighted {
		return s.LabelHighlight
	}
	return s.Label
}

func fileLabels(s Styles, cursorFile chess.File, fileOrder []chess.File) string {
	var sb strings.Builder
	sb.WriteString("  ") // align with rank label width
	for _, f := range fileOrder {
		sb.WriteString(s.labelStyle(f == cursorFile).Render(fmt.Sprintf(" %c ", rune('a')+rune(f))))
	}
	return sb.String()
}

func renderSquare(sq chess.Square, piece chess.Piece, opts RenderOptions, s Styles) string {
	f, r := sq.File(), sq.Rank()
	isLight := (int(f)+int(r))%2 == 1
	isCursor := sq == opts.Cursor
	isSelected := opts.Selected != nil && sq == *opts.Selected
	isValidDest := opts.ValidDests[sq]
	isLastMove := (opts.LastMoveFrom != nil && sq == *opts.LastMoveFrom) ||
		(opts.LastMoveTo != nil && sq == *opts.LastMoveTo)

	var bg lipgloss.Style
	switch {
	case isCursor && (isValidDest || isSelected) && isLight:
		bg = s.SelectedLight
	case isCursor && (isValidDest || isSelected):
		bg = s.SelectedDark
	case isCursor && isLight:
		bg = s.CursorLight
	case isCursor:
		bg = s.CursorDark
	case isSelected && isLight:
		bg = s.SelectedLight
	case isSelected:
		bg = s.SelectedDark
	case isValidDest && isLight:
		bg = s.CursorLight
	case isValidDest:
		bg = s.CursorDark
	case isLastMove && isLight:
		bg = s.LastMoveLight
	case isLastMove:
		bg = s.LastMoveDark
	case isLight:
		bg = s.LightSquare
	default:
		bg = s.DarkSquare
	}

	if piece == chess.NoPiece {
		return bg.Render("   ")
	}

	fg := lipgloss.Color("#1a1a1a")
	if piece.Color() == chess.White {
		fg = "#FFFFFF"
	}

	return bg.Foreground(fg).Bold(true).Render(" " + pieces[piece.Type()] + " ")
}
