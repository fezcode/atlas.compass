package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fezcode/atlas.compass/internal/tui"
)

var Version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Printf("atlas.compass v%s\n", Version)
		return
	}

	p := tea.NewProgram(tui.NewMainModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running atlas.compass: %v\n", err)
		os.Exit(1)
	}
}