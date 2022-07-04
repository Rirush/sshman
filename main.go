package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	prog := tea.NewProgram(NewConnectionsModel())
	if err := prog.Start(); err != nil {
		fmt.Printf("Program failed: %v\n", err)
	}
}
