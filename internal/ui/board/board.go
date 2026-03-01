package board

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/notnil/chess"
)

type Styles struct {
	LightSquare  lipgloss.Style
	DarkSquare   lipgloss.Style
	CursorLight  lipgloss.Style
	CursorDark   lipgloss.Style
	Label        lipgloss.Style
}

func NewStyles(renderer *lipgloss.Renderer) Styles {
	return Styles{
		LightSquare: renderer.NewStyle().Background(lipgloss.Color("#F0D9B5")),
		DarkSquare:  renderer.NewStyle().Background(lipgloss.Color("#B58863")),
		CursorLight: renderer.NewStyle().Background(lipgloss.Color("#CCD277")),
		CursorDark:  renderer.NewStyle().Background(lipgloss.Color("#AAA23A")),
		Label: renderer.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#9b9b9b", Dark: "#5c5c5c"}),
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

// Render draws the board from White's perspective.
func Render(pos *chess.Position, s Styles, cursor chess.Square) string {
	b := pos.Board()
	var sb strings.Builder

	sb.WriteString(fileLabels(s))
	sb.WriteString("\n")

	for r := chess.Rank8; r >= chess.Rank1; r-- {
		sb.WriteString(s.Label.Render(fmt.Sprintf("%d ", int(r)+1)))
		for f := chess.FileA; f <= chess.FileH; f++ {
			sq := chess.NewSquare(f, r)
			sb.WriteString(renderSquare(sq, b.Piece(sq), cursor, s))
		}
		sb.WriteString(s.Label.Render(fmt.Sprintf(" %d", int(r)+1)))
		sb.WriteString("\n")
	}

	sb.WriteString(fileLabels(s))
	return sb.String()
}

func fileLabels(s Styles) string {
	var sb strings.Builder
	sb.WriteString("  ") // align with rank label width
	for f := chess.FileA; f <= chess.FileH; f++ {
		sb.WriteString(s.Label.Render(fmt.Sprintf(" %c ", rune('a')+rune(f))))
	}
	return sb.String()
}

func renderSquare(sq chess.Square, piece chess.Piece, cursor chess.Square, s Styles) string {
	f, r := sq.File(), sq.Rank()
	isLight := (int(f)+int(r))%2 == 1

	var bg lipgloss.Style
	switch {
	case sq == cursor && isLight:
		bg = s.CursorLight
	case sq == cursor:
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
