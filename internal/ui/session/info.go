package session

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
)

type TerminalInfo struct {
	Term    string
	Profile string
	Width   int
	Height  int
	BG      string
}

func ExtractTerminalInfo(renderer *lipgloss.Renderer, pty ssh.Pty) TerminalInfo {
	bg := "light"
	if renderer.HasDarkBackground() {
		bg = "dark"
	}

	return TerminalInfo{
		Term:    pty.Term,
		Profile: renderer.ColorProfile().Name(),
		Width:   pty.Window.Width,
		Height:  pty.Window.Height,
		BG:      bg,
	}
}
