# DNS Switcher ⚡

A high-performance DNS management tool that provides a **Native Windows GUI** and a sleek **Terminal UI (TUI)** for Unix-based systems. Optimize your internet speed and privacy with one click.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%11%20macOS-blue)
![Architecture](https://img.shields.io/badge/Interface-GUI%20%26%20TUI-brightgreen)

## ✨ Features

- **Dual-Interface System**:
    - **Windows**: Beautiful native GUI inspired by V2RayNG (built with Fyne).
    - **Linux/macOS**: Professional, interactive Terminal UI (built with BubbleTea).
- **Embedded Branded Icon**: Windows executable comes with a custom blue shield icon.
- **Speed Sorting**: One-click "Sort by Speed" to instantly find the lowest latency servers.
- **Live Monitoring**: Real-time dashboard showing uptime, latency, and query success/failure.
- **Safe & Secure**: 
    - Automatic Windows UAC prompt for Administrator elevation.
    - Automatic DNS configuration backup (Linux).
- **Dedicated Disconnect**: Easily revert to system default DNS settings with a single button/hotkey.
- **Extensive Provider List**: Over 15+ pre-configured high-performance DNS servers.

## 🚀 Installation

### Windows (GUI)
1. Download `dns-switcher.exe` from the latest release.
2. Run as Administrator.
3. Enjoy the V2RayNG-inspired dark theme.

### Linux/macOS (TUI)
1. Download the binary for your platform.
2. Run with sudo: `sudo ./dns-switcher`

---

## 🛠 Building from Source

### Prerequisites
- Go 1.21+
- **Windows**: [MSYS2](https://www.msys64.org/) (for GCC/CGO) is required for the GUI.

### Windows Build (GUI)
```powershell
$env:CGO_ENABLED="1"; $env:Path += ";C:\msys64\mingw64\bin"; & "C:\Program Files\Go\bin\go.exe" build -ldflags="-s -w -H windowsgui" -o dns-switcher.exe .
```

### Linux/macOS Build (TUI)
```bash
go build -o dns-switcher .
```

---

## 🧭 Navigation

### Windows (GUI)
- **Sidebar**: Toggle between DNS Servers, Monitoring, and Settings.
- **Connect**: Click a provider's card to switch immediately.
- **Sort**: Use the "Sort by Speed" button in the header.

### Linux/macOS (TUI)
- `↑/↓` or `j/k`: Navigate through providers.
- `Enter`: Select a provider.
- `r`: Refresh latency in monitor mode.
- `c`: Change DNS (go back).
- `q`: Quit.

---

## 📋 Supported DNS Providers
- **Privacy**: Shecan, AdGuard, CleanBrowsing.
- **Performance**: Cloudflare, Google, OpenDNS, Quad9.
- **Regional**: Radar, Electro, Begzar, 403.
- **Custom**: "Add Custom DNS" allows you to paste any IP (e.g., `8.8.8.8 1.1.1.1`).

## ⚙️ How It Works
- **Windows**: Uses PowerShell `Set-DnsClientServerAddress` and `ipconfig /flushdns`.
- **Linux**: Manages `/etc/resolv.conf` and restarts `systemd-resolved`.
- **macOS**: Uses the system `networksetup` utility for active services.

## 📄 License
MIT License - see [LICENSE](LICENSE) file for details.

## 🤝 Acknowledgments
- [Fyne](https://fyne.io/) - Native Windows GUI framework.
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework.
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Terminal styling.

Created by [Harry](https://github.com/Har2yQn78)
