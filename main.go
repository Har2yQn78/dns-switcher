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

	fmt.Println("Current DNS servers:")
	currentDNS, err := GetCurrentDNS()
	if err != nil {
		fmt.Printf("  (Could not read: %v)\n", err)
	} else if len(currentDNS) == 0 {
		fmt.Println("  (none)")
	} else {
		for _, dns := range currentDNS {
			fmt.Printf("  • %s\n", dns)
		}
	}
	fmt.Println()

	// Start the bubbletea program
	p := tea.NewProgram(initialModel())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	m := finalModel.(model)
	if m.selected >= 0 && !m.quitting {
		provider := providers[m.selected]
		fmt.Printf("\n✓ Selected: %s\n", provider.Name)

		fmt.Println("\nUpdating DNS settings...")
		if err := UpdateResolvConf(provider); err != nil {
			fmt.Printf("✗ Error: %v\n", err)
			os.Exit(1)
		}

		if err := RestartSystemdResolved(); err != nil {
			fmt.Printf("Warning: %v\n", err)
		}

		fmt.Println("\n✓ DNS updated successfully!")

		fmt.Println("\nNew DNS servers:")
		newDNS, _ := GetCurrentDNS()
		for _, dns := range newDNS {
			fmt.Printf("  • %s\n", dns)
		}
	}
}
