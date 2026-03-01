package session

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	wishtea "github.com/charmbracelet/wish/bubbletea"

	"terminalchess/internal/ui/root"
	"terminalchess/internal/ui/styles"
)

func TeaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// activeterm middleware guarantees this is present.
	pty, _, _ := s.Pty()

	renderer := wishtea.MakeRenderer(s)

	m := root.New(root.Props{
		Width:  pty.Window.Width,
		Height: pty.Window.Height,
		Styles: styles.New(renderer),
	})

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}
