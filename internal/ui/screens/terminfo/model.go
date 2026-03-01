package terminfo

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"terminalchess/internal/ui/navigate"
)

type Props struct {
	Term      string
	Profile   string
	Width     int
	Height    int
	BG        string
	TxtStyle  lipgloss.Style
	QuitStyle lipgloss.Style
}

type Model struct {
	term      string
	profile   string
	width     int
	height    int
	bg        string
	txtStyle  lipgloss.Style
	quitStyle lipgloss.Style
}

func NewModel(p Props) Model {
	return Model{
		term:      p.Term,
		profile:   p.Profile,
		width:     p.Width,
		height:    p.Height,
		bg:        p.BG,
		txtStyle:  p.TxtStyle,
		quitStyle: p.QuitStyle,
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc", "b":
			return m, navigate.To(navigate.Menu)
		}
	}
	return m, nil
}

func (m Model) View() string {
	s := fmt.Sprintf(
		"Your term is %s\nYour window size is %dx%d\nBackground: %s\nColor Profile: %s",
		m.term, m.width, m.height, m.bg, m.profile,
	)
	return m.txtStyle.Render(s) + "\n\n" + m.quitStyle.Render("Press 'esc' to go back  •  'q' to quit\n")
}
