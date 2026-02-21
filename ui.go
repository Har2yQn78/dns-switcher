//go:build !windows

package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func RunApp() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

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
	cursor       int
	selected     int
	quitting     bool
	inputMode    bool
	customInput  string
	customError  string
	monitorMode  bool
	monitorStats MonitorStats
	scrollOffset int
	termHeight   int
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

func (m model) visibleRows() int {
	if m.termHeight <= 0 {
		return len(providers)
	}
	rows := m.termHeight - 12
	if rows < 5 {
		rows = 5
	}
	if rows > len(providers) {
		rows = len(providers)
	}
	return rows
}

func (m model) adjustScroll() model {
	visible := m.visibleRows()
	if m.cursor < m.scrollOffset {
		m.scrollOffset = m.cursor
	}
	if m.cursor >= m.scrollOffset+visible {
		m.scrollOffset = m.cursor - visible + 1
	}
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}
	return m
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

	case tea.WindowSizeMsg:
		m.termHeight = msg.Height
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
				m = m.adjustScroll()
			}

		case "down", "j":
			if m.cursor < len(providers)-1 {
				m.cursor++
				m = m.adjustScroll()
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
			b.WriteString(errorStyle.Render("  "+m.customError) + "\n\n")
		}

		b.WriteString(helpStyle.Render("  Example: 8.8.8.8,1.1.1.1 or 8.8.8.8 1.1.1.1") + "\n")
		b.WriteString(helpStyle.Render("  enter: confirm • esc: cancel") + "\n")

		return b.String()
	}

	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(titleStyle.Render("  DNS Changer") + "\n")
	b.WriteString(helpStyle.Render("  Press q or ctrl+c to quit") + "\n\n")

	// Column content widths (characters of visible text)
	const nameWidth = 20
	const serverWidth = 40
	const latWidth = 8

	topBorder := fmt.Sprintf("  ┌%s┬%s┬%s┐",
		strings.Repeat("─", nameWidth+2),
		strings.Repeat("─", serverWidth+2),
		strings.Repeat("─", latWidth+2))
	b.WriteString(borderStyle.Render(topBorder) + "\n")

	hProvider := fmt.Sprintf("%-*s", nameWidth, "Provider")
	hServers := fmt.Sprintf("%-*s", serverWidth, "DNS Servers")
	hLatency := fmt.Sprintf("%-*s", latWidth, "Latency")
	header := fmt.Sprintf("  │ %s │ %s │ %s │",
		headerStyle.Render(hProvider),
		headerStyle.Render(hServers),
		headerStyle.Render(hLatency))
	b.WriteString(header + "\n")

	sep := fmt.Sprintf("  ├%s┼%s┼%s┤",
		strings.Repeat("─", nameWidth+2),
		strings.Repeat("─", serverWidth+2),
		strings.Repeat("─", latWidth+2))
	b.WriteString(borderStyle.Render(sep) + "\n")

	// Viewport scrolling
	visible := m.visibleRows()
	startIdx := m.scrollOffset
	endIdx := startIdx + visible
	if endIdx > len(providers) {
		endIdx = len(providers)
	}

	for i := startIdx; i < endIdx; i++ {
		provider := providers[i]

		providerName := provider.Name
		if m.cursor == i {
			providerName = "▸ " + providerName
		} else {
			providerName = "  " + providerName
		}
		paddedName := fmt.Sprintf("%-*s", nameWidth, providerName)

		var servers string
		if len(provider.Servers) == 0 {
			servers = ""
		} else if len(provider.Servers) == 1 {
			servers = provider.Servers[0]
		} else {
			servers = fmt.Sprintf("%-17s %s", provider.Servers[0], provider.Servers[1])
		}
		paddedServers := fmt.Sprintf("%-*s", serverWidth, servers)

		var latencyStr string
		var latStyle lipgloss.Style

		if provider.Latency == -1 {
			latencyStr = "N/A"
			latStyle = failedLatencyStyle
		} else if provider.Latency < 20 {
			latencyStr = fmt.Sprintf("%dms", provider.Latency)
			latStyle = fastLatencyStyle
		} else if provider.Latency < 50 {
			latencyStr = fmt.Sprintf("%dms", provider.Latency)
			latStyle = mediumLatencyStyle
		} else {
			latencyStr = fmt.Sprintf("%dms", provider.Latency)
			latStyle = slowLatencyStyle
		}
		paddedLatency := fmt.Sprintf("%*s", latWidth, latencyStr)

		if m.cursor == i {
			row := fmt.Sprintf("  │ %s │ %s │ %s │",
				selectedRowStyle.Render(paddedName),
				selectedRowStyle.Render(paddedServers),
				selectedRowStyle.Render(paddedLatency))
			b.WriteString(row + "\n")
		} else {
			row := fmt.Sprintf("  │ %s │ %s │ %s │",
				normalRowStyle.Render(paddedName),
				serverStyle.Render(paddedServers),
				latStyle.Render(paddedLatency))
			b.WriteString(row + "\n")
		}
	}

	btmBorder := fmt.Sprintf("  └%s┴%s┴%s┘",
		strings.Repeat("─", nameWidth+2),
		strings.Repeat("─", serverWidth+2),
		strings.Repeat("─", latWidth+2))
	b.WriteString(borderStyle.Render(btmBorder) + "\n")

	// Scroll indicators
	var scrollInfo []string
	if startIdx > 0 {
		scrollInfo = append(scrollInfo, fmt.Sprintf("▲ %d more above", startIdx))
	}
	if endIdx < len(providers) {
		scrollInfo = append(scrollInfo, fmt.Sprintf("▼ %d more below", len(providers)-endIdx))
	}
	if len(scrollInfo) > 0 {
		b.WriteString(helpStyle.Render("  "+strings.Join(scrollInfo, " • ")) + "\n")
	}
	b.WriteString("\n")

	help := helpStyle.Render("  Use ↑/↓ or j/k to navigate • enter to select • q to quit")
	b.WriteString(help + "\n")

	return b.String()
}
