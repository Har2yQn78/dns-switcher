package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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
		return "Goodbye!\n"
	}

	s := "DNS Changer - Select a provider:\n\n"

	for i, provider := range providers {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		servers := strings.Join(provider.Servers, ", ")

		s += fmt.Sprintf("%s %-20s %s\n", cursor, provider.Name, servers)
	}

	s += "\nUse ↑/↓ arrows to move, Enter to select, q to quit\n"

	return s
}
