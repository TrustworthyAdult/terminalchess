package main

import (
	"os"

	"github.com/charmbracelet/log"
	"terminalchess/internal/server"
)

func main() {
	cfg := server.Config{
		Host:        "localhost",
		Port:        "23234",
		HostKeyPath: ".ssh/id_ed25519",
	}

	s, err := server.New(cfg)
	if err != nil {
		log.Error("Could not start server", "error", err)
		os.Exit(1)
	}

	if err := server.RunUntilSignal(s); err != nil {
		log.Error("Server exited with error", "error", err)
		os.Exit(1)
	}
}
