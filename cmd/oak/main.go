package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jam-computing/oak/pkg/components"
)

func main() {
    model := components.GetModel()

    fmt.Println("Hello, WorlD!")
    p := tea.NewProgram(model)
    if _, err := p.Run(); err != nil {
        fmt.Println("Error running program", err)
        os.Exit(1)
    }
}
