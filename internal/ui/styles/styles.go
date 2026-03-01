package styles

import (
	"github.com/charmbracelet/lipgloss"

	"terminalchess/internal/ui/board"
)

var (
	green  = lipgloss.AdaptiveColor{Light: "#4a7c59", Dark: "#769656"}
	cream  = lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#ffffd7"}
	selFg  = lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#ffffd7"}
	subtle = lipgloss.AdaptiveColor{Light: "#9b9b9b", Dark: "#5c5c5c"}
)

type Styles struct {
	Panel        lipgloss.Style
	Title        lipgloss.Style
	Hint         lipgloss.Style
	Body         lipgloss.Style
	Cursor       lipgloss.Style
	NormalItem   lipgloss.Style
	SelectedItem lipgloss.Style
	Board        board.Styles
}

func New(renderer *lipgloss.Renderer) Styles {
	return Styles{
		Panel: renderer.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(green).
			Padding(1, 4),
		Title: renderer.NewStyle().
			Foreground(green).
			Bold(true),
		Hint: renderer.NewStyle().
			Foreground(subtle),
		Body: renderer.NewStyle().
			Foreground(cream),
		Cursor: renderer.NewStyle().
			Foreground(green),
		NormalItem: renderer.NewStyle().
			Foreground(subtle).
			Padding(0, 1),
		SelectedItem: renderer.NewStyle().
			Background(green).
			Foreground(selFg).
			Bold(true).
			Padding(0, 1),
		Board: board.NewStyles(renderer),
	}
}
