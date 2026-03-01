package menu

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"terminalchess/internal/ui/navigate"
	"terminalchess/internal/ui/styles"
)

// itemWidth is the fixed width of each menu item, ensuring the selected
// background highlight fills a consistent row.
const itemWidth = 18

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Quit  key.Binding
}

var keys = keyMap{
	Up:    key.NewBinding(key.WithKeys("up", "k")),
	Down:  key.NewBinding(key.WithKeys("down", "j")),
	Enter: key.NewBinding(key.WithKeys("enter")),
	Quit:  key.NewBinding(key.WithKeys("ctrl+c", "q")),
}

type menuItem struct {
	title string
	dest  navigate.Screen
	quit  bool
}

type Props struct {
	Styles styles.Styles
}

type Model struct {
	items  []menuItem
	cursor int
	styles styles.Styles
}

func NewModel(p Props) Model {
	return Model{
		items: []menuItem{
			{title: "Play Game", dest: navigate.Game},
			{title: "Quit", quit: true},
		},
		styles: p.Styles,
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(k, keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(k, keys.Down):
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case key.Matches(k, keys.Enter):
			it := m.items[m.cursor]
			if it.quit {
				return m, tea.Quit
			}
			return m, navigate.To(it.dest)
		case key.Matches(k, keys.Quit):
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	s := m.styles

	title := s.Title.Render("♟  Terminal Chess")

	var rows []string
	for i, it := range m.items {
		var line string
		if i == m.cursor {
			line = s.Cursor.Render("▸") + " " + s.SelectedItem.Width(itemWidth).Render(it.title)
		} else {
			line = "  " + s.NormalItem.Width(itemWidth).Render(it.title)
		}
		rows = append(rows, line)
	}

	hint := s.Hint.Render("↑/↓ move   enter select")

	body := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		strings.Join(rows, "\n"),
		"",
		hint,
	)

	return s.Panel.Render(body)
}
