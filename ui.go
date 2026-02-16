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

	// Latency color styles
	fastLatencyStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#50FA7B")) // Green < 20ms

	mediumLatencyStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#F1FA8C")) // Yellow 20-50ms

	slowLatencyStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF5555")) // Red > 50ms

	failedLatencyStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#6272A4")) // Gray (failed/not tested)
)

type model struct {
	cursor       int
	selected     int
	quitting     bool
	inputMode    bool
	customInput  string
	customError  string
}

func initialModel() model {
	return model{
		cursor:      0,
		selected:    -1,
		quitting:    false,
		inputMode:   false,
		customInput: "",
		customError: "",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.inputMode {
			switch msg.String() {
			case "ctrl+c":
				m.quitting = true
				return m, tea.Quit

			case "esc":
				m.inputMode = false
				m.customInput = ""
				m.customError = ""

			case "enter":
				if m.customInput == "" {
					m.customError = "Please enter at least one DNS server"
				} else {
					servers := parseCustomDNS(m.customInput)
					if len(servers) == 0 {
						m.customError = "Invalid DNS format. Use: 8.8.8.8,1.1.1.1"
					} else {
						customProvider := DNSProvider{
							Name:    "Custom DNS",
							Servers: servers,
							Latency: -1,
						}
						providers = append(providers, customProvider)
						m.selected = len(providers) - 1
						return m, tea.Quit
					}
				}

			case "backspace":
				if len(m.customInput) > 0 {
					m.customInput = m.customInput[:len(m.customInput)-1]
					m.customError = ""
				}

			default:
				if len(msg.String()) == 1 {
					m.customInput += msg.String()
					m.customError = ""
				}
			}
			return m, nil
		}

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
			if providers[m.cursor].Name == "Add Custom DNS" {
				m.inputMode = true
				m.customInput = ""
				m.customError = ""
			} else {
				m.selected = m.cursor
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func parseCustomDNS(input string) []string {
	input = strings.ReplaceAll(input, ",", " ")

	parts := strings.Fields(input)

	var servers []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			servers = append(servers, part)
		}
	}

	return servers
}

func (m model) View() string {
	if m.quitting {
		return titleStyle.Render("Goodbye!\n")
	}

	if m.inputMode {
		var b strings.Builder

		b.WriteString("\n")
		b.WriteString(titleStyle.Render("  Add Custom DNS") + "\n\n")

		b.WriteString(infoStyle.Render("  Enter DNS servers (comma or space separated):") + "\n\n")
		b.WriteString(fmt.Sprintf("  > %s_\n\n", m.customInput))

		if m.customError != "" {
			b.WriteString(errorStyle.Render("  " + m.customError) + "\n\n")
		}

		b.WriteString(helpStyle.Render("  Example: 8.8.8.8,1.1.1.1 or 8.8.8.8 1.1.1.1") + "\n")
		b.WriteString(helpStyle.Render("  enter: confirm • esc: cancel") + "\n")

		return b.String()
	}

	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(titleStyle.Render("  DNS Changer") + "\n")
	b.WriteString(helpStyle.Render("  Press q or ctrl+c to quit") + "\n\n")
	border := borderStyle.Render("  ┌──────────────────────┬──────────────────────────────────────────┬──────────┐")
	b.WriteString(border + "\n")

	header := fmt.Sprintf("  │ %-20s │ %-40s │ %-8s │",
		headerStyle.Render("Provider"),
		headerStyle.Render("DNS Servers"),
		headerStyle.Render("Latency"),
	)
	b.WriteString(header + "\n")

	separator := borderStyle.Render("  ├──────────────────────┼──────────────────────────────────────────┼──────────┤")
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

		var servers string
		if len(provider.Servers) == 0 {
			servers = ""
		} else if len(provider.Servers) == 1 {
			servers = fmt.Sprintf("%-16s", provider.Servers[0])
		} else {
			servers = fmt.Sprintf("%-16s %-16s", provider.Servers[0], provider.Servers[1])
		}

		servers = fmt.Sprintf("%-40s", servers)
		var latencyStr string
		var latencyStyle lipgloss.Style

		if provider.Latency == -1 {
			latencyStr = "     N/A"
			latencyStyle = failedLatencyStyle
		} else if provider.Latency < 20 {
			latencyStr = fmt.Sprintf("%6dms", provider.Latency)
			latencyStyle = fastLatencyStyle
		} else if provider.Latency < 50 {
			latencyStr = fmt.Sprintf("%6dms", provider.Latency)
			latencyStyle = mediumLatencyStyle
		} else {
			latencyStr = fmt.Sprintf("%6dms", provider.Latency)
			latencyStyle = slowLatencyStyle
		}

		if m.cursor == i {
			row := fmt.Sprintf("  │ %-20s │ %-40s │ %s │",
				rowStyle.Render(providerName),
				rowStyle.Render(servers),
				rowStyle.Render(latencyStr),
			)
			b.WriteString(row + "\n")
		} else {
			row := fmt.Sprintf("  │ %-20s │ %-40s │ %s │",
				rowStyle.Render(providerName),
				serverStyle.Render(servers),
				latencyStyle.Render(latencyStr),
			)
			b.WriteString(row + "\n")
		}
	}
	bottomBorder := borderStyle.Render("  └──────────────────────┴────────────────────────────────────┴──────────┘")
	b.WriteString(bottomBorder + "\n\n")

	help := helpStyle.Render("  Use ↑/↓ or j/k to navigate • enter to select • q to quit")
	b.WriteString(help + "\n")

	return b.String()
}
