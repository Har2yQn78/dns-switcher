package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	boxStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272A4"))

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F2"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555"))

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BD93F9")).
			Bold(true)
)

func printBox(title string, content []string) {
	width := 60

	fmt.Println(boxStyle.Render("  ┌" + repeatStr("─", width-4) + "┐"))
	fmt.Println(boxStyle.Render("  │ ") + labelStyle.Render(title) + boxStyle.Render(repeatStr(" ", width-len(title)-6) + "│"))
	fmt.Println(boxStyle.Render("  ├" + repeatStr("─", width-4) + "┤"))

	for _, line := range content {
		padding := width - len(stripANSI(line)) - 6
		if padding < 0 {
			padding = 0
		}
		fmt.Println(boxStyle.Render("  │ ") + line + boxStyle.Render(repeatStr(" ", padding) + "│"))
	}

	fmt.Println(boxStyle.Render("  └" + repeatStr("─", width-4) + "┘"))
	fmt.Println()
}

func repeatStr(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

func stripANSI(str string) string {
	inEscape := false
	count := 0
	for _, r := range str {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		count++
	}
	result := ""
	for i := 0; i < count; i++ {
		result += " "
	}
	return result
}

func main() {
	if os.Geteuid() != 0 {
		fmt.Println(errorStyle.Render("Error: Please run this program as root (use sudo)"))
		os.Exit(1)
	}

	currentDNS, err := GetCurrentDNS()
	var dnsLines []string
	if err != nil {
		dnsLines = append(dnsLines, infoStyle.Render("Could not read current DNS"))
	} else if len(currentDNS) == 0 {
		dnsLines = append(dnsLines, infoStyle.Render("No DNS servers configured"))
	} else {
		for _, dns := range currentDNS {
			dnsLines = append(dnsLines, infoStyle.Render(dns))
		}
	}
	printBox("Current DNS Servers", dnsLines)

	p := tea.NewProgram(initialModel())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf(errorStyle.Render("Error: %v\n"), err)
		os.Exit(1)
	}

	m := finalModel.(model)
	if m.selected >= 0 && !m.quitting {
		provider := providers[m.selected]

		fmt.Println(labelStyle.Render("  Selected: ") + infoStyle.Render(provider.Name))
		fmt.Println()

		statusLines := []string{}
		statusLines = append(statusLines, infoStyle.Render("Updating DNS configuration..."))

		if err := UpdateResolvConf(provider); err != nil {
			fmt.Println(errorStyle.Render("  Error: " + err.Error()))
			os.Exit(1)
		}
		statusLines = append(statusLines, successStyle.Render("Configuration updated"))

		if err := RestartSystemdResolved(); err != nil {
			statusLines = append(statusLines, errorStyle.Render("Warning: "+err.Error()))
		} else {
			statusLines = append(statusLines, successStyle.Render("Service restarted"))
		}

		printBox("Update Status", statusLines)

		newDNS, _ := GetCurrentDNS()
		var newDNSLines []string
		for _, dns := range newDNS {
			newDNSLines = append(newDNSLines, infoStyle.Render(dns))
		}
		printBox("New DNS Servers", newDNSLines)

		if provider.Name != "Reset to Default" {
			validationLines := []string{}
			validationLines = append(validationLines, infoStyle.Render("Testing DNS resolution..."))

			success, validationErr := ValidateDNS(provider.Servers)
			if success {
				validationLines = append(validationLines, successStyle.Render("All DNS servers responding"))
			} else {
				validationLines = append(validationLines, errorStyle.Render("Validation failed: "+validationErr.Error()))
			}

			printBox("DNS Validation", validationLines)
		}
	}
}
