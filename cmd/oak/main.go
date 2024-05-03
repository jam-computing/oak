package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jam-computing/oak/pkg/components"
)

func main() {
    m := components.GetModel()

    p := tea.NewProgram(m)
    if _, err := p.Run(); err != nil {
        fmt.Println("Error running program", err)
        os.Exit(1)
    }
}
