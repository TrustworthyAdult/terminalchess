package menu

import (
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"terminalchess/internal/ui/navigate"
	"terminalchess/internal/ui/styles"
	"terminalchess/internal/ui/tui"
)

// itemWidth is the fixed width of each menu item, ensuring the selected
// background highlight fills a consistent row.
const itemWidth = 24

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Quit  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.Enter, k.Quit}}
}

var keys = keyMap{
	Up:    key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	Down:  key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	Enter: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	Quit:  key.NewBinding(key.WithKeys("ctrl+c", "q"), key.WithHelp("q", "quit")),
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
	help   help.Model
	offset float64
}

func NewModel(p Props) Model {
	return Model{
		items: []menuItem{
			{title: "Local Game", dest: navigate.Game},
			{title: "Play Against Computer", dest: navigate.Setup},
			{title: "Quit", quit: true},
		},
		styles: p.Styles,
		help:   help.New(),
	}
}

func (m Model) Init() tea.Cmd { return tui.Tick() }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tui.TickMsg:
		m.offset = math.Mod(m.offset+2, 360)
		return m, tui.Tick()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, keys.Down):
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case key.Matches(msg, keys.Enter):
			it := m.items[m.cursor]
			if it.quit {
				return m, tea.Quit
			}
			return m, navigate.To(it.dest)
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	s := m.styles

	// Greyscale overrides — the rainbow border provides all the colour.
	title := s.Title.Foreground(lipgloss.Color("#FFFFFF")).Render("♟  Terminal Chess")
	cursorStyle := s.Cursor.Foreground(lipgloss.Color("#AAAAAA"))
	selectedStyle := s.SelectedItem.
		Background(lipgloss.Color("#505050")).
		Foreground(lipgloss.Color("#FFFFFF"))

	var rows []string
	for i, it := range m.items {
		var line string
		if i == m.cursor {
			line = cursorStyle.Render("▸") + " " + selectedStyle.Width(itemWidth).Render(it.title)
		} else {
			line = "  " + s.NormalItem.Width(itemWidth).Render(it.title)
		}
		rows = append(rows, line)
	}

	body := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		strings.Join(rows, "\n"),
		"",
		m.help.View(keys),
	)

	inner := lipgloss.NewStyle().Padding(1, 4).Render(body)
	return tui.RainbowBorder(inner, m.offset)
}
