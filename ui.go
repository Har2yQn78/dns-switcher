package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time

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

	fastLatencyStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#50FA7B"))

	mediumLatencyStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#F1FA8C"))

	slowLatencyStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF5555"))

	failedLatencyStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#6272A4"))
)

type model struct {
	cursor         int
	selected       int
	quitting       bool
	inputMode      bool
	customInput    string
	customError    string
	monitorMode    bool
	monitorStats   MonitorStats
}

type MonitorStats struct {
	ProviderName   string
	CurrentDNS     []string
	QueriesSuccess int
	QueriesFailed  int
	LastLatency    int
	Uptime         int
}

func initialModel() model {
	return model{
		cursor:       0,
		selected:     -1,
		quitting:     false,
		inputMode:    false,
		customInput:  "",
		customError:  "",
		monitorMode:  false,
		monitorStats: MonitorStats{},
	}
}

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	if m.monitorMode {
		return doTick()
	}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if m.monitorMode {
			m.monitorStats.Uptime++

			if len(m.monitorStats.CurrentDNS) > 0 {
				latency := TestDNSLatency(m.monitorStats.CurrentDNS[0])
				if latency > 0 {
					m.monitorStats.LastLatency = latency
					m.monitorStats.QueriesSuccess++
				} else {
					m.monitorStats.QueriesFailed++
				}
			}

			return m, doTick()
		}
		return m, nil

	case tea.KeyMsg:
		if m.monitorMode {
			switch msg.String() {
			case "ctrl+c", "q":
				m.quitting = true
				return m, tea.Quit
			case "r":
				if len(m.monitorStats.CurrentDNS) > 0 {
					latency := TestDNSLatency(m.monitorStats.CurrentDNS[0])
					if latency > 0 {
						m.monitorStats.LastLatency = latency
						m.monitorStats.QueriesSuccess++
					} else {
						m.monitorStats.QueriesFailed++
					}
				}
			case "c":
				m.monitorMode = false
				m.selected = -1
				return m, tea.Quit
			}
			return m, nil
		}

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

func formatDuration(seconds int) string {
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}
	minutes := seconds / 60
	secs := seconds % 60
	if minutes < 60 {
		return fmt.Sprintf("%dm %ds", minutes, secs)
	}
	hours := minutes / 60
	mins := minutes % 60
	return fmt.Sprintf("%dh %dm", hours, mins)
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

	if m.monitorMode {
		var b strings.Builder

		b.WriteString("\n")
		b.WriteString(titleStyle.Render("  DNS Monitoring Dashboard") + "\n\n")

		b.WriteString(fmt.Sprintf("  %s %s\n\n",
			headerStyle.Render("Provider:"),
			infoStyle.Render(m.monitorStats.ProviderName)))

		b.WriteString(headerStyle.Render("  Current DNS Servers:") + "\n")
		for _, dns := range m.monitorStats.CurrentDNS {
			b.WriteString(fmt.Sprintf("    • %s\n", serverStyle.Render(dns)))
		}
		b.WriteString("\n")

		border := borderStyle.Render("  ┌────────────────────────┬──────────────┐")
		b.WriteString(border + "\n")

		uptimeStr := formatDuration(m.monitorStats.Uptime)
		b.WriteString(fmt.Sprintf("  │ %-22s │ %-12s │\n",
			headerStyle.Render("Uptime"),
			infoStyle.Render(uptimeStr)))

		var latencyStr string
		var latencyColor lipgloss.Style
		if m.monitorStats.LastLatency < 20 {
			latencyStr = fmt.Sprintf("%dms", m.monitorStats.LastLatency)
			latencyColor = fastLatencyStyle
		} else if m.monitorStats.LastLatency < 50 {
			latencyStr = fmt.Sprintf("%dms", m.monitorStats.LastLatency)
			latencyColor = mediumLatencyStyle
		} else if m.monitorStats.LastLatency > 0 {
			latencyStr = fmt.Sprintf("%dms", m.monitorStats.LastLatency)
			latencyColor = slowLatencyStyle
		} else {
			latencyStr = "N/A"
			latencyColor = failedLatencyStyle
		}
		b.WriteString(fmt.Sprintf("  │ %-22s │ %-12s │\n",
			headerStyle.Render("Current Latency"),
			latencyColor.Render(latencyStr)))

		b.WriteString(fmt.Sprintf("  │ %-22s │ %-12s │\n",
			headerStyle.Render("Queries Success"),
			fastLatencyStyle.Render(fmt.Sprintf("%d", m.monitorStats.QueriesSuccess))))

		if m.monitorStats.QueriesFailed > 0 {
			b.WriteString(fmt.Sprintf("  │ %-22s │ %-12s │\n",
				headerStyle.Render("Queries Failed"),
				slowLatencyStyle.Render(fmt.Sprintf("%d", m.monitorStats.QueriesFailed))))
		} else {
			b.WriteString(fmt.Sprintf("  │ %-22s │ %-12s │\n",
				headerStyle.Render("Queries Failed"),
				infoStyle.Render("0")))
		}

		bottomBorder := borderStyle.Render("  └────────────────────────┴──────────────┘")
		b.WriteString(bottomBorder + "\n\n")

		b.WriteString(helpStyle.Render("  r: refresh • c: change DNS • q: quit") + "\n")

		return b.String()
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

	bottomBorder := borderStyle.Render("  └──────────────────────┴──────────────────────────────────────────┴──────────┘")
	b.WriteString(bottomBorder + "\n\n")

	help := helpStyle.Render("  Use ↑/↓ or j/k to navigate • enter to select • q to quit")
	b.WriteString(help + "\n")

	return b.String()
}
