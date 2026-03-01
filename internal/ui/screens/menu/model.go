package menu

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"terminalchess/internal/ui/navigate"
)

type item struct {
	title string
	dest  navigate.Screen
	quit  bool
}

func (i item) FilterValue() string { return i.title }
func (i item) Title() string       { return i.title }
func (i item) Description() string { return "" }

type Props struct {
	Width     int
	Height    int
	TxtStyle  lipgloss.Style
	QuitStyle lipgloss.Style
}

type Model struct {
	list list.Model
}

func NewModel(p Props) Model {
	items := []list.Item{
		item{title: "Play Game", dest: navigate.Game},
		item{title: "Quit", quit: true},
	}

	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, p.Width, p.Height)
	l.Title = "Terminal Chess"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return Model{list: l}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
		if selected, ok := m.list.SelectedItem().(item); ok {
			if selected.quit {
				return m, tea.Quit
			}
			return m, navigate.To(selected.dest)
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.list.View()
}
