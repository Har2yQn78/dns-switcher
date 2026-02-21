//go:build windows

package main

import (
	"fmt"
	"image/color"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func RunApp() {
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

type AppState struct {
	mu              sync.Mutex
	activeProvider  string
	activeDNS       []string
	connected       bool
	monitorUptime   int
	monitorSuccess  int
	monitorFailed   int
	lastLatency     int
	monitorRunning  bool
	monitorStop     chan struct{}
}

var state = &AppState{}

func makeSidebar(contentArea *fyne.Container, w fyne.Window) fyne.CanvasObject {
	logo := canvas.NewText("⚡ DNS Switcher", colorPrimary)
	logo.TextSize = 20
	logo.TextStyle = fyne.TextStyle{Bold: true}
	logo.Alignment = fyne.TextAlignCenter

	version := canvas.NewText("v2.0", colorTextSecondary)
	version.TextSize = 11
	version.Alignment = fyne.TextAlignCenter

	headerBox := container.NewVBox(
		container.NewPadded(logo),
		version,
		widget.NewSeparator(),
	)

	btnServers := widget.NewButtonWithIcon("DNS Servers", theme.ComputerIcon(), func() {
		contentArea.Objects = []fyne.CanvasObject{makeServersPanel(contentArea, w)}
		contentArea.Refresh()
	})
	btnMonitor := widget.NewButtonWithIcon("Monitor", theme.InfoIcon(), func() {
		contentArea.Objects = []fyne.CanvasObject{makeMonitorPanel()}
		contentArea.Refresh()
	})
	btnSettings := widget.NewButtonWithIcon("Settings", theme.SettingsIcon(), func() {
		contentArea.Objects = []fyne.CanvasObject{makeSettingsPanel()}
		contentArea.Refresh()
	})

	btnDisconnect := widget.NewButtonWithIcon("Disconnect", theme.ContentClearIcon(), func() {
		state.mu.Lock()
		state.connected = false
		state.activeProvider = ""
		state.activeDNS = nil
		state.mu.Unlock()
		stopMonitor()
		
		var resetProv DNSProvider
		for _, p := range providers {
			if p.Name == "Reset to Default" {
				resetProv = p
				break
			}
		}
		if resetProv.Name != "" {
			UpdateResolvConf(resetProv)
		}
		
		dialog.ShowInformation("Disconnected", "DNS has been reset to system defaults.", w)
		contentArea.Objects = []fyne.CanvasObject{makeServersPanel(contentArea, w)}
		contentArea.Refresh()
	})

	for _, btn := range []*widget.Button{btnServers, btnMonitor, btnSettings, btnDisconnect} {
		btn.Importance = widget.LowImportance
	}
	btnDisconnect.Importance = widget.LowImportance

	navBox := container.NewVBox(
		btnServers,
		btnMonitor,
		btnSettings,
		widget.NewSeparator(),
		btnDisconnect,
	)

	quitBtn := widget.NewButtonWithIcon("Quit", theme.LogoutIcon(), func() {
		stopMonitor()
		fyne.CurrentApp().Quit()
	})
	quitBtn.Importance = widget.DangerImportance

	sidebar := container.NewBorder(
		headerBox,
		container.NewPadded(quitBtn),
		nil, nil,
		container.NewPadded(navBox),
	)

	bg := canvas.NewRectangle(colorSidebarBg)
	return container.NewStack(bg, container.NewPadded(sidebar))
}

func makeStatusBar() fyne.CanvasObject {
	statusDot := canvas.NewCircle(colorDisconnected)
	statusDot.StrokeWidth = 0
	statusDotSized := container.NewStack(
		container.New(layout.NewGridWrapLayout(fyne.NewSize(10, 10)), statusDot),
	)

	statusLabel := canvas.NewText("  Disconnected", colorTextSecondary)
	statusLabel.TextSize = 12

	currentDNS, err := GetCurrentDNS()
	if err == nil && len(currentDNS) > 0 {
		statusDot.FillColor = colorConnected
		statusLabel.Text = fmt.Sprintf("  Active: %s", currentDNS[0])
		statusLabel.Color = colorTextPrimary
	}

	state.mu.Lock()
	if state.connected {
		statusDot.FillColor = colorConnected
		statusLabel.Text = fmt.Sprintf("  Connected: %s", state.activeProvider)
		statusLabel.Color = colorSuccess
	}
	state.mu.Unlock()

	adapterLabel := canvas.NewText("", colorTextSecondary)
	adapterLabel.TextSize = 11

	statusRow := container.NewHBox(statusDotSized, statusLabel, layout.NewSpacer(), adapterLabel)
	bg := canvas.NewRectangle(colorSurface)
	return container.NewStack(bg, container.NewPadded(statusRow))
}

func makeServersPanel(contentArea *fyne.Container, w fyne.Window) fyne.CanvasObject {
	title := canvas.NewText("DNS Servers", colorTextPrimary)
	title.TextSize = 22
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("Select a DNS provider to connect", colorTextSecondary)
	subtitle.TextSize = 13

	sortBtn := widget.NewButtonWithIcon("Sort by Speed", theme.ListIcon(), func() {
		SortProvidersByLatency()
		contentArea.Objects = []fyne.CanvasObject{makeServersPanel(contentArea, w)}
		contentArea.Refresh()
	})
	sortBtn.Importance = widget.LowImportance

	header := container.NewVBox(
		container.NewBorder(nil, nil, nil, sortBtn, title),
		subtitle,
		widget.NewSeparator(),
	)

	cards := container.NewVBox()
	for i, p := range providers {
		idx := i
		prov := p
		card := makeProviderCard(idx, prov, w)
		cards.Add(card)
	}

	scrollable := container.NewVScroll(cards)
	scrollable.SetMinSize(fyne.NewSize(600, 400))

	return container.NewBorder(
		container.NewPadded(container.NewPadded(header)),
		nil, nil, nil,
		container.NewPadded(scrollable),
	)
}

func makeProviderCard(idx int, prov DNSProvider, w fyne.Window) fyne.CanvasObject {
	name := canvas.NewText(prov.Name, colorTextPrimary)
	name.TextSize = 15
	name.TextStyle = fyne.TextStyle{Bold: true}

	var serversStr string
	if len(prov.Servers) == 0 {
		serversStr = "Custom..."
	} else {
		for i, s := range prov.Servers {
			if i > 0 {
				serversStr += "  •  "
			}
			serversStr += s
		}
	}
	servers := canvas.NewText(serversStr, colorTextSecondary)
	servers.TextSize = 12

	var latencyWidget fyne.CanvasObject
	if prov.Latency == -1 {
		badge := canvas.NewText("  N/A  ", colorTextSecondary)
		badge.TextSize = 11
		latencyWidget = badge
	} else {
		latStr := fmt.Sprintf("  %dms  ", prov.Latency)
		badgeColor := colorSuccess
		if prov.Latency >= 50 {
			badgeColor = colorError
		} else if prov.Latency >= 20 {
			badgeColor = colorWarning
		}
		badge := canvas.NewText(latStr, badgeColor)
		badge.TextSize = 11
		badge.TextStyle = fyne.TextStyle{Bold: true}
		latencyWidget = badge
	}

	info := container.NewVBox(name, servers)

	var actionBtn *widget.Button
	if prov.Name == "Add Custom DNS" {
		actionBtn = widget.NewButtonWithIcon("Add", theme.ContentAddIcon(), func() {
			showCustomDNSDialog(w)
		})
	} else {
		actionBtn = widget.NewButtonWithIcon("Connect", theme.NavigateNextIcon(), func() {
			connectToProvider(idx, w)
		})
		actionBtn.Importance = widget.HighImportance
	}

	rightSide := container.NewHBox(latencyWidget, actionBtn)
	row := container.NewBorder(nil, nil, nil, rightSide, info)

	bg := canvas.NewRectangle(colorSurface)
	bg.CornerRadius = 10

	padded := container.NewPadded(row)
	card := container.NewStack(bg, container.NewPadded(padded))

	return container.NewPadded(card)
}

func connectToProvider(idx int, w fyne.Window) {
	provider := providers[idx]

	prog := dialog.NewProgressInfinite("Connecting...",
		fmt.Sprintf("Setting DNS to %s", provider.Name), w)
	prog.Show()

	go func() {
		err := UpdateResolvConf(provider)
		prog.Hide()

		if err != nil {
			dialog.ShowError(fmt.Errorf("Failed to set DNS: %v", err), w)
			return
		}

		_ = RestartSystemdResolved()

		state.mu.Lock()
		state.activeProvider = provider.Name
		state.activeDNS = provider.Servers
		state.connected = true
		state.mu.Unlock()

		if provider.Name != "Reset to Default" && len(provider.Servers) > 0 {
			success, valErr := ValidateDNS(provider.Servers)
			if success {
				dialog.ShowInformation("Connected",
					fmt.Sprintf("✅ Successfully connected to %s\nDNS servers are responding.", provider.Name), w)
			} else {
				dialog.ShowInformation("Warning",
					fmt.Sprintf("⚠️ DNS set to %s but validation issue:\n%v", provider.Name, valErr), w)
			}
		} else {
			dialog.ShowInformation("Reset",
				"DNS has been reset to default settings.", w)
			state.mu.Lock()
			state.connected = false
			state.mu.Unlock()
		}

		startMonitor()
	}()
}

func showCustomDNSDialog(w fyne.Window) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("e.g. 8.8.8.8, 1.1.1.1")

	items := []*widget.FormItem{
		widget.NewFormItem("DNS Servers", entry),
	}

	dlg := dialog.NewForm("Add Custom DNS", "Connect", "Cancel", items, func(ok bool) {
		if !ok || entry.Text == "" {
			return
		}
		servers := parseCustomDNS(entry.Text)
		if len(servers) == 0 {
			dialog.ShowError(fmt.Errorf("Invalid DNS format.\nUse: 8.8.8.8, 1.1.1.1"), w)
			return
		}
		customProvider := DNSProvider{
			Name:    "Custom DNS",
			Servers: servers,
			Latency: -1,
		}
		providers = append(providers, customProvider)
		connectToProvider(len(providers)-1, w)
	}, w)
	dlg.Resize(fyne.NewSize(400, 200))
	dlg.Show()
}

func makeMonitorPanel() fyne.CanvasObject {
	title := canvas.NewText("Monitoring Dashboard", colorTextPrimary)
	title.TextSize = 22
	title.TextStyle = fyne.TextStyle{Bold: true}

	state.mu.Lock()
	provName := state.activeProvider
	connected := state.connected
	dns := state.activeDNS
	uptime := state.monitorUptime
	success := state.monitorSuccess
	failed := state.monitorFailed
	latency := state.lastLatency
	state.mu.Unlock()

	if !connected {
		noConn := canvas.NewText("No active DNS connection. Select a provider first.", colorTextSecondary)
		noConn.TextSize = 14
		return container.NewVBox(
			container.NewPadded(title),
			widget.NewSeparator(),
			container.NewCenter(container.NewPadded(noConn)),
		)
	}

	provLabel := canvas.NewText(fmt.Sprintf("Provider: %s", provName), colorPrimary)
	provLabel.TextSize = 16
	provLabel.TextStyle = fyne.TextStyle{Bold: true}

	var dnsStr string
	for i, d := range dns {
		if i > 0 {
			dnsStr += "  •  "
		}
		dnsStr += d
	}
	dnsLabel := canvas.NewText(fmt.Sprintf("DNS: %s", dnsStr), colorTextSecondary)
	dnsLabel.TextSize = 13

	uptimeCard := makeStatCard("⏱ Uptime", formatDuration(uptime), colorPrimary)
	latencyCard := makeStatCard("📡 Latency", formatLatency(latency), latencyColor(latency))
	successCard := makeStatCard("✅ Success", fmt.Sprintf("%d", success), colorSuccess)
	failedCard := makeStatCard("❌ Failed", fmt.Sprintf("%d", failed), colorError)

	statsGrid := container.NewGridWithColumns(4,
		uptimeCard, latencyCard, successCard, failedCard,
	)

	refreshBtn := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() {
		state.mu.Lock()
		if len(state.activeDNS) > 0 {
			lat := TestDNSLatency(state.activeDNS[0])
			if lat > 0 {
				state.lastLatency = lat
				state.monitorSuccess++
			} else {
				state.monitorFailed++
			}
		}
		state.mu.Unlock()
	})

	return container.NewVBox(
		container.NewPadded(title),
		widget.NewSeparator(),
		container.NewPadded(provLabel),
		container.NewPadded(dnsLabel),
		widget.NewSeparator(),
		container.NewPadded(statsGrid),
		container.NewPadded(container.NewCenter(refreshBtn)),
	)
}

func makeStatCard(label string, value string, col color.Color) fyne.CanvasObject {
	bg := canvas.NewRectangle(colorSurface)
	bg.CornerRadius = 10

	lbl := canvas.NewText(label, colorTextSecondary)
	lbl.TextSize = 12
	lbl.Alignment = fyne.TextAlignCenter

	val := canvas.NewText(value, col)
	val.TextSize = 22
	val.TextStyle = fyne.TextStyle{Bold: true}
	val.Alignment = fyne.TextAlignCenter

	content := container.NewVBox(
		container.NewCenter(lbl),
		container.NewCenter(val),
	)

	return container.NewStack(bg, container.NewPadded(content))
}

func formatLatency(ms int) string {
	if ms <= 0 {
		return "N/A"
	}
	return fmt.Sprintf("%dms", ms)
}

func latencyColor(ms int) color.Color {
	if ms <= 0 {
		return colorTextSecondary
	} else if ms < 20 {
		return colorSuccess
	} else if ms < 50 {
		return colorWarning
	}
	return colorError
}

func startMonitor() {
	stopMonitor()
	state.mu.Lock()
	state.monitorRunning = true
	state.monitorUptime = 0
	state.monitorSuccess = 0
	state.monitorFailed = 0
	state.monitorStop = make(chan struct{})
	stopCh := state.monitorStop
	state.mu.Unlock()

	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-stopCh:
				return
			case <-ticker.C:
				state.mu.Lock()
				state.monitorUptime += 2
				if len(state.activeDNS) > 0 {
					lat := TestDNSLatency(state.activeDNS[0])
					if lat > 0 {
						state.lastLatency = lat
						state.monitorSuccess++
					} else {
						state.monitorFailed++
					}
				}
				state.mu.Unlock()
			}
		}
	}()
}

func stopMonitor() {
	state.mu.Lock()
	defer state.mu.Unlock()
	if state.monitorRunning && state.monitorStop != nil {
		close(state.monitorStop)
		state.monitorRunning = false
	}
}

func makeSettingsPanel() fyne.CanvasObject {
	title := canvas.NewText("Settings", colorTextPrimary)
	title.TextSize = 22
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("Application preferences", colorTextSecondary)
	subtitle.TextSize = 13

	retestBtn := widget.NewButtonWithIcon("Re-test All Latencies", theme.ViewRefreshIcon(), func() {
		go func() {
			TestAllProviders()
		}()
	})

	aboutTitle := canvas.NewText("About", colorPrimary)
	aboutTitle.TextSize = 16
	aboutTitle.TextStyle = fyne.TextStyle{Bold: true}

	aboutText := canvas.NewText("DNS Switcher v2.0", colorTextSecondary)
	aboutText.TextSize = 12

	authorText := canvas.NewText("Created by Harry", colorTextSecondary)
	authorText.TextSize = 12

	aboutBg := canvas.NewRectangle(colorSurface)
	aboutBg.CornerRadius = 10

	aboutCard := container.NewStack(aboutBg, container.NewPadded(container.NewVBox(
		aboutTitle,
		aboutText,
		authorText,
	)))

	return container.NewVBox(
		container.NewPadded(title),
		container.NewPadded(subtitle),
		widget.NewSeparator(),
		container.NewPadded(retestBtn),
		widget.NewSeparator(),
		container.NewPadded(aboutCard),
	)
}
