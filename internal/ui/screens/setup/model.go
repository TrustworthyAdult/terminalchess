package setup

import (
	"math"
	"math/rand"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/notnil/chess"

	"terminalchess/internal/ui/navigate"
	"terminalchess/internal/ui/screens/game"
	"terminalchess/internal/ui/styles"
	"terminalchess/internal/ui/tui"
)

type section int

const (
	colorSection      section = iota
	difficultySection
	startSection
)

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Enter key.Binding
	Back  key.Binding
	Quit  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Left, k.Right, k.Up, k.Back}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.Left, k.Right, k.Enter, k.Back, k.Quit}}
}

var keys = keyMap{
	Up:    key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	Down:  key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	Left:  key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/→", "choose")),
	Right: key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("↑/↓", "row")),
	Enter: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "start")),
	Back:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	Quit:  key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
}

type Props struct {
	Styles styles.Styles
}

type Model struct {
	styles     styles.Styles
	section    section
	color      int // 0=White 1=Black 2=Random
	difficulty int // 0=Easy  1=Medium 2=Hard
	help       help.Model
	offset     float64
}

func NewModel(p Props) Model {
	return Model{
		styles:     p.Styles,
		difficulty: 1, // default Medium
		help:       help.New(),
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
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Back):
			return m, navigate.To(navigate.Menu)
		case key.Matches(msg, keys.Up):
			m.section = (m.section + 2) % 3
		case key.Matches(msg, keys.Down):
			m.section = (m.section + 1) % 3
		case key.Matches(msg, keys.Left):
			switch m.section {
			case colorSection:
				m.color = (m.color + 2) % 3
			case difficultySection:
				m.difficulty = (m.difficulty + 2) % 3
			case startSection:
				m.section = difficultySection
			}
		case key.Matches(msg, keys.Right):
			switch m.section {
			case colorSection:
				m.color = (m.color + 1) % 3
			case difficultySection:
				m.difficulty = (m.difficulty + 1) % 3
			case startSection:
				// no-op on start row
			}
		case key.Matches(msg, keys.Enter):
			switch m.section {
			case colorSection:
				m.section = difficultySection
			case difficultySection:
				m.section = startSection
			case startSection:
				return m, m.startGame()
			}
		}
	}
	return m, nil
}

func (m Model) startGame() tea.Cmd {
	// Determine the human's chosen color first, then assign the opposite to the computer.
	humanColor := chess.White
	switch m.color {
	case 1:
		humanColor = chess.Black
	case 2:
		if rand.Intn(2) == 1 {
			humanColor = chess.Black
		}
	}
	computerColor := chess.Black
	if humanColor == chess.Black {
		computerColor = chess.White
	}
	skillLevel := []int{5, 12, 20}[m.difficulty]
	return navigate.ToWithConfig(navigate.Game, game.Config{
		ComputerColor: &computerColor,
		SkillLevel:    skillLevel,
	})
}

func (m Model) View() string {
	chosen := lipgloss.NewStyle().Bold(true).
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#FFFFFF"))
	unchosen := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888"))

	focusedLabel := lipgloss.NewStyle().
		Bold(true).Foreground(lipgloss.Color("#FFFFFF"))
	unfocusedLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666"))

	renderRow := func(sec section, label string, opts []string, selected int) string {
		lStyle := unfocusedLabel
		if m.section == sec {
			lStyle = focusedLabel
		}
		var choices []string
		for i, opt := range opts {
			padded := " " + opt + " "
			if i == selected {
				choices = append(choices, chosen.Render(padded))
			} else {
				choices = append(choices, unchosen.Render(padded))
			}
		}
		return lStyle.Render(label) + "  " + strings.Join(choices, "  ")
	}

	colorRow := renderRow(colorSection, "Color:     ",
		[]string{"White", "Black", "Random"}, m.color)
	diffRow := renderRow(difficultySection, "Difficulty:",
		[]string{"Easy", "Medium", "Hard"}, m.difficulty)

	startStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	if m.section == startSection {
		startStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF"))
	}
	startRow := startStyle.Render("          ▸ Start Game")

	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF")).
		Render("♟  Play vs Computer")

	body := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		colorRow,
		diffRow,
		"",
		startRow,
		"",
		m.help.View(keys),
	)

	inner := lipgloss.NewStyle().Padding(1, 4).Render(body)
	return tui.RainbowBorder(inner, m.offset)
}
