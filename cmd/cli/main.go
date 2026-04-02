package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"torque/cmd/cli/tui"
)

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	godotenv.Load()

	cfg := tui.Config{
		KratosAdminURL: getenv("KRATOS_ADMIN_URL", "http://localhost:4434"),
		KetoReadURL:    getenv("KETO_READ_URL", "http://localhost:4466"),
		KetoWriteURL:   getenv("KETO_WRITE_URL", "http://localhost:4467"),
	}

	p := tea.NewProgram(tui.New(cfg))
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
