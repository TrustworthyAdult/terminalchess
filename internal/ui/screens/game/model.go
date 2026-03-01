package game

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"terminalchess/internal/ui/navigate"
	"terminalchess/internal/ui/styles"
)

type Props struct {
	Styles styles.Styles
}

type Model struct {
	styles styles.Styles
}

func NewModel(p Props) Model {
	return Model{styles: p.Styles}
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
	s := m.styles

	content := lipgloss.JoinVertical(lipgloss.Center,
		s.Body.Render("Chess game coming soon..."),
		"",
		s.Hint.Render("esc  go back"),
	)

	return s.Panel.Render(content)
}
