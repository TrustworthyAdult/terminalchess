package menu

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"terminalchess/internal/ui/navigate"
	"terminalchess/internal/ui/styles"
)

// itemWidth is the fixed width of each menu item, ensuring the selected
// background highlight fills a consistent row.
const itemWidth = 18

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

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
			{title: "Play Game", dest: navigate.Game},
			{title: "Quit", quit: true},
		},
		styles: p.Styles,
		help:   help.New(),
	}
}

func (m Model) Init() tea.Cmd { return tick() }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.offset = math.Mod(m.offset+2, 360)
		return m, tick()
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
	return rainbowBorder(inner, m.offset)
}

// rainbowBorder wraps content in a rounded border where each character's hue
// is offset by its clockwise position around the perimeter, creating a
// spinning gradient effect as offset advances.
func rainbowBorder(content string, offset float64) string {
	lines := strings.Split(strings.TrimRight(content, "\n"), "\n")
	maxW := 0
	for _, l := range lines {
		if w := lipgloss.Width(l); w > maxW {
			maxW = w
		}
	}
	height := len(lines)
	perimeter := 2*maxW + 2*height + 4

	colorAt := func(pos int) lipgloss.Color {
		hue := math.Mod(offset+float64(pos)/float64(perimeter)*360, 360)
		return lipgloss.Color(hslToHex(hue, 1.0, 0.6))
	}
	ch := func(pos int, r string) string {
		return lipgloss.NewStyle().Foreground(colorAt(pos)).Render(r)
	}

	var sb strings.Builder

	// Top edge (positions 0 … maxW+1)
	sb.WriteString(ch(0, "╭"))
	for i := 0; i < maxW; i++ {
		sb.WriteString(ch(1+i, "─"))
	}
	sb.WriteString(ch(maxW+1, "╮"))
	sb.WriteString("\n")

	// Content rows: right border travels top→bottom (maxW+2 … maxW+1+height),
	// left border travels bottom→top (2*maxW+height+4 … 2*maxW+2*height+3).
	for i, line := range lines {
		rightPos := maxW + 2 + i
		leftPos := 2*maxW + 2*height + 3 - i
		sb.WriteString(ch(leftPos, "│"))
		sb.WriteString(line)
		if lw := lipgloss.Width(line); lw < maxW {
			sb.WriteString(strings.Repeat(" ", maxW-lw))
		}
		sb.WriteString(ch(rightPos, "│"))
		sb.WriteString("\n")
	}

	// Bottom edge: ╰ at 2*maxW+height+3, then right→left (2*maxW+height+2 … maxW+height+3), ╯ at maxW+height+2
	sb.WriteString(ch(2*maxW+height+3, "╰"))
	for i := 0; i < maxW; i++ {
		sb.WriteString(ch(2*maxW+height+2-i, "─"))
	}
	sb.WriteString(ch(maxW+height+2, "╯"))

	return sb.String()
}

// hslToHex converts HSL (h∈[0,360], s∈[0,1], l∈[0,1]) to a CSS hex colour.
func hslToHex(h, s, l float64) string {
	h /= 360
	var r, g, b float64
	if s == 0 {
		r, g, b = l, l, l
	} else {
		hue2rgb := func(p, q, t float64) float64 {
			if t < 0 {
				t += 1
			}
			if t > 1 {
				t -= 1
			}
			switch {
			case t < 1.0/6.0:
				return p + (q-p)*6*t
			case t < 0.5:
				return q
			case t < 2.0/3.0:
				return p + (q-p)*(2.0/3.0-t)*6
			}
			return p
		}
		q := l * (1 + s)
		if l >= 0.5 {
			q = l + s - l*s
		}
		p := 2*l - q
		r = hue2rgb(p, q, h+1.0/3.0)
		g = hue2rgb(p, q, h)
		b = hue2rgb(p, q, h-1.0/3.0)
	}
	return fmt.Sprintf("#%02X%02X%02X", int(r*255+0.5), int(g*255+0.5), int(b*255+0.5))
}
