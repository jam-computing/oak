package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jam-computing/oak/pkg/components"
)

func main() {
    m := components.GetModel(components.ViewingAnimations)

    p := tea.NewProgram(m, tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Println("Error running program", err)
        os.Exit(1)
    }
}
