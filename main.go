package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
)

func TestAllProviders() {
	for i := range providers {
		if providers[i].Name == "Reset to Default" || providers[i].Name == "Add Custom DNS" {
			providers[i].Latency = -1
			continue
		}

		if len(providers[i].Servers) > 0 {
			latency := TestDNSLatency(providers[i].Servers[0])
			providers[i].Latency = latency
		}
	}
}

func SortProvidersByLatency() {
	var special []DNSProvider
	var normal []DNSProvider

	for _, p := range providers {
		if p.Name == "Reset to Default" || p.Name == "Add Custom DNS" {
			special = append(special, p)
		} else {
			normal = append(normal, p)
		}
	}

	sort.Slice(normal, func(i, j int) bool {
		if normal[i].Latency == -1 {
			return false
		}
		if normal[j].Latency == -1 {
			return true
		}
		return normal[i].Latency < normal[j].Latency
	})

	providers = append(normal, special...)
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

func main() {
	a := app.New()
	a.Settings().SetTheme(newDNSTheme())

	w := a.NewWindow("DNS Switcher")
	w.Resize(fyne.NewSize(920, 620))
	w.SetFixedSize(false)
	w.CenterOnScreen()

	if runtime.GOOS == "windows" && !IsAdmin() {
		verb := "runas"
		exe, _ := os.Executable()
		cwd, _ := os.Getwd()
		args := strings.Join(os.Args[1:], " ")

		cmd := exec.Command("powershell", "Start-Process", fmt.Sprintf("'%s'", exe), "-ArgumentList", fmt.Sprintf("'%s'", args), "-Verb", verb, "-WorkingDirectory", fmt.Sprintf("'%s'", cwd))
		err := cmd.Start()
		if err == nil {
			os.Exit(0)
		}
	}

	loadingLabel := canvas.NewText("Testing DNS providers...", colorTextPrimary)
	loadingLabel.TextSize = 16
	loadingLabel.Alignment = fyne.TextAlignCenter
	w.SetContent(container.NewCenter(loadingLabel))
	w.Show()

	go func() {
		TestAllProviders()

		contentArea := container.NewStack()
		contentArea.Objects = []fyne.CanvasObject{makeServersPanel(contentArea, w)}

		sidebar := makeSidebar(contentArea, w)
		sidebarSized := container.New(
			layout.NewGridWrapLayout(fyne.NewSize(220, 0)),
			sidebar,
		)

		divider := canvas.NewRectangle(colorDivider)
		divider.SetMinSize(fyne.NewSize(1, 0))

		mainBody := container.NewBorder(
			nil, nil,
			container.NewHBox(sidebarSized, divider),
			nil,
			contentArea,
		)

		statusBar := makeStatusBar()

		fullLayout := container.NewBorder(
			nil,
			statusBar,
			nil, nil,
			mainBody,
		)

		w.SetContent(fullLayout)
		w.Canvas().Refresh(fullLayout)
	}()

	if !IsAdmin() {
		dialog.ShowError(fmt.Errorf(
			"Administrator privileges required.\n\n"+
				"Windows: Right-click → Run as Administrator\n"+
				"Linux/macOS: sudo ./dns-switcher"), w)
	}

	a.Run()
}
