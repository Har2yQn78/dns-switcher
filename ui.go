package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF79C6")).
			Bold(true)

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BD93F9")).
			Bold(true)

	selectedRowStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#44475A")).
				Foreground(lipgloss.Color("#50FA7B")).
				Bold(true)

	normalRowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F2"))

	serverStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8BE9FD"))

	borderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272A4"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272A4")).
			Italic(true)
)

type model struct {
	cursor   int
	selected int
	quitting bool
}

func initialModel() model {
	return model{
		cursor:   0,
		selected: -1,
		quitting: false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(providers)-1 {
				m.cursor++
			}

		case "enter", " ":
			m.selected = m.cursor
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return titleStyle.Render("Goodbye!\n")
	}

	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(titleStyle.Render("  DNS Changer") + "\n")
	b.WriteString(helpStyle.Render("  Press q or ctrl+c to quit") + "\n\n")
	border := borderStyle.Render("  ┌──────────────────────┬────────────────────────────────────┐")
	b.WriteString(border + "\n")

	header := fmt.Sprintf("  │ %-20s │ %-34s │",
		headerStyle.Render("Provider"),
		headerStyle.Render("DNS Servers"),
	)
	b.WriteString(header + "\n")

	separator := borderStyle.Render("  ├──────────────────────┼────────────────────────────────────┤")
	b.WriteString(separator + "\n")

	for i, provider := range providers {
		var rowStyle lipgloss.Style
		if m.cursor == i {
			rowStyle = selectedRowStyle
		} else {
			rowStyle = normalRowStyle
		}

		providerName := provider.Name
		if m.cursor == i {
			providerName = "> " + providerName
		} else {
			providerName = "  " + providerName
		}

		var serverList []string
		for _, server := range provider.Servers {
			serverList = append(serverList, fmt.Sprintf("%-15s", server))
		}
		servers := strings.Join(serverList, " ")

		if len(servers) > 34 {
			servers = servers[:31] + "..."
		}

		if m.cursor == i {
			row := fmt.Sprintf("  │ %-20s │ %-34s │",
				rowStyle.Render(providerName),
				rowStyle.Render(servers),
			)
			b.WriteString(row + "\n")
		} else {
			row := fmt.Sprintf("  │ %-20s │ %-34s │",
				rowStyle.Render(providerName),
				serverStyle.Render(servers),
			)
			b.WriteString(row + "\n")
		}
	}

	bottomBorder := borderStyle.Render("  └──────────────────────┴────────────────────────────────────┘")
	b.WriteString(bottomBorder + "\n\n")

	help := helpStyle.Render("  Use ↑/↓ or j/k to navigate • enter to select • q to quit")
	b.WriteString(help + "\n")

	return b.String()
}
