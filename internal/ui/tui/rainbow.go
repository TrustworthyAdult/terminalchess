package tui

import (
	"fmt"
	"math"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TickMsg time.Time

func Tick() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg { return TickMsg(t) })
}

// RainbowBorder wraps content in a rounded border where each character's hue
// is offset by its clockwise position around the perimeter, creating a
// spinning gradient effect as offset advances.
func RainbowBorder(content string, offset float64) string {
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
