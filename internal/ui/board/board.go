package board

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/notnil/chess"
)

type Styles struct {
	LightSquare   lipgloss.Style
	DarkSquare    lipgloss.Style
	CursorLight   lipgloss.Style
	CursorDark    lipgloss.Style
	SelectedLight lipgloss.Style
	SelectedDark  lipgloss.Style
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
		Label: renderer.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#9b9b9b", Dark: "#5c5c5c"}),
		LabelHighlight: renderer.NewStyle().
			Foreground(cursorLight).
			Bold(true),
	}
}

var whitePieces = map[chess.PieceType]string{
	chess.King:   "♔",
	chess.Queen:  "♕",
	chess.Rook:   "♖",
	chess.Bishop: "♗",
	chess.Knight: "♘",
	chess.Pawn:   "♙",
}

var blackPieces = map[chess.PieceType]string{
	chess.King:   "♚",
	chess.Queen:  "♛",
	chess.Rook:   "♜",
	chess.Bishop: "♝",
	chess.Knight: "♞",
	chess.Pawn:   "♙",
}

type RenderOptions struct {
	Cursor     chess.Square
	Selected   *chess.Square
	ValidDests map[chess.Square]bool
}

// Render draws the board from White's perspective.
func Render(pos *chess.Position, s Styles, opts RenderOptions) string {
	b := pos.Board()
	cursorFile := opts.Cursor.File()
	cursorRank := opts.Cursor.Rank()
	var sb strings.Builder

	sb.WriteString(fileLabels(s, cursorFile))
	sb.WriteString("\n")

	for r := chess.Rank8; r >= chess.Rank1; r-- {
		rankLabel := s.labelStyle(r == cursorRank)
		sb.WriteString(rankLabel.Render(fmt.Sprintf("%d ", int(r)+1)))
		for f := chess.FileA; f <= chess.FileH; f++ {
			sq := chess.NewSquare(f, r)
			sb.WriteString(renderSquare(sq, b.Piece(sq), opts, s))
		}
		sb.WriteString(rankLabel.Render(fmt.Sprintf(" %d", int(r)+1)))
		sb.WriteString("\n")
	}

	sb.WriteString(fileLabels(s, cursorFile))
	return sb.String()
}

func (s Styles) labelStyle(highlighted bool) lipgloss.Style {
	if highlighted {
		return s.LabelHighlight
	}
	return s.Label
}

func fileLabels(s Styles, cursorFile chess.File) string {
	var sb strings.Builder
	sb.WriteString("  ") // align with rank label width
	for f := chess.FileA; f <= chess.FileH; f++ {
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
	case isLight:
		bg = s.LightSquare
	default:
		bg = s.DarkSquare
	}

	if piece == chess.NoPiece {
		return bg.Render("   ")
	}

	var (
		symbol string
		fg     lipgloss.Color
	)

	if piece.Color() == chess.White {
		symbol = whitePieces[piece.Type()]
		fg = "#FFFFFF"
	} else {
		symbol = blackPieces[piece.Type()]
		fg = "#1a1a1a"
	}

	return bg.Foreground(fg).Bold(true).Render(" " + symbol + " ")
}
