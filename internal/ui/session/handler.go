package session

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	wishtea "github.com/charmbracelet/wish/bubbletea"

	"terminalchess/internal/ui/screens/terminfo"
)

func TeaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// activeterm middleware guarantees this is present.
	pty, _, _ := s.Pty()

	renderer := wishtea.MakeRenderer(s)

	info := ExtractTerminalInfo(renderer, pty)

	styles := NewStyles(renderer)

	m := terminfo.NewModel(
		terminfo.Props{
			Term:      info.Term,
			Profile:   info.Profile,
			Width:     info.Width,
			Height:    info.Height,
			BG:        info.BG,
			TxtStyle:  styles.Txt,
			QuitStyle: styles.Quit,
		},
	)

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}
