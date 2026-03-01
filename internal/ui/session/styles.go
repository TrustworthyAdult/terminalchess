package session

import (
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	Txt  lipgloss.Style
	Quit lipgloss.Style
}

func NewStyles(renderer *lipgloss.Renderer) Styles {
	return Styles{
		Txt:  renderer.NewStyle().Foreground(lipgloss.Color("10")),
		Quit: renderer.NewStyle().Foreground(lipgloss.Color("8")),
	}
}
