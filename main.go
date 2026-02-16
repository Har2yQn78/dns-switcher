package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if os.Geteuid() != 0 {
		fmt.Println("Error: Please run this program as root (use sudo)")
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	m := finalModel.(model)
	if m.selected >= 0 && !m.quitting {
		provider := providers[m.selected]
		fmt.Printf("\n✓ Selected: %s (%s)\n", provider.Name, provider.Servers)
	}
}
