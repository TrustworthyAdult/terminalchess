package server

import (
	"net"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	wishtea "github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"

	"terminalchess/internal/ui/session"
)

type Config struct {
	Host        string
	Port        string
	HostKeyPath string
}

func New(cfg Config) (*ssh.Server, error) {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(cfg.Host, cfg.Port)),
		wish.WithHostKeyPath(cfg.HostKeyPath),
		wish.WithMiddleware(
			wishtea.Middleware(func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
				return session.TeaHandler(s)
			}),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		return nil, err
	}

	log.Info("Starting SSH server", "host", cfg.Host, "port", cfg.Port)
	return s, nil
}
