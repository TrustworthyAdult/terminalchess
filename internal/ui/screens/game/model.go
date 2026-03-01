package game

import (
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"

	"terminalchess/internal/ui/navigate"
)

type Props struct {
	TxtStyle  lipgloss.Style
	QuitStyle lipgloss.Style
}

type Model struct {
	txtStyle  lipgloss.Style
	quitStyle lipgloss.Style
}

func NewModel(p Props) Model {
	return Model{txtStyle: p.TxtStyle, quitStyle: p.QuitStyle}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "q", "esc":
			return m, navigate.To(navigate.Menu)
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	return m.txtStyle.Render("Chess game coming soon...") +
		"\n\n" + m.quitStyle.Render("Press 'esc' to return to menu")
}
