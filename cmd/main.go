package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fezcode/atlas.compass/internal/tui"
)

func main() {
	p := tea.NewProgram(tui.NewMainModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running atlas.compass: %v\n", err)
		os.Exit(1)
	}
}