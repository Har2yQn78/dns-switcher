# DNS Switcher ⚡

A high-performance DNS management tool that provides a **Native Windows GUI** and a sleek **Terminal UI (TUI)** for Unix-based systems. Optimize your internet speed and privacy with one click.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-blue)
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

### Install from Release (Recommended)

Download the latest release from the [releases page](https://github.com/Har2yQn78/dns-switcher/releases):

**Linux:**

```bash
# Download the binary
wget https://github.com/Har2yQn78/dns-switcher/releases/latest/download/dns-switcher-linux-amd64

# Install it
sudo mv dns-switcher-linux-amd64 /usr/local/bin/dns-switcher
sudo chmod +x /usr/local/bin/dns-switcher
```

**macOS (Intel):**

```bash
# Download the binary
curl -L https://github.com/Har2yQn78/dns-switcher/releases/latest/download/dns-switcher-macos-amd64 -o dns-switcher

# Install it
sudo mv dns-switcher /usr/local/bin/
sudo chmod +x /usr/local/bin/dns-switcher
```

**macOS (Apple Silicon):**

```bash
# Download the binary
curl -L https://github.com/Har2yQn78/dns-switcher/releases/latest/download/dns-switcher-macos-arm64 -o dns-switcher

# Install it
sudo mv dns-switcher /usr/local/bin/
sudo chmod +x /usr/local/bin/dns-switcher
```

**Windows (64-bit):**

```powershell
# Download using PowerShell
Invoke-WebRequest -Uri "https://github.com/Har2yQn78/dns-switcher/releases/latest/download/dns-switcher-windows-amd64.exe" -OutFile "dns-switcher.exe"

# Move to a directory in your PATH (optional)
# Or run directly from current directory
```

**Windows (32-bit):**

```powershell
# Download using PowerShell
Invoke-WebRequest -Uri "https://github.com/Har2yQn78/dns-switcher/releases/latest/download/dns-switcher-windows-386.exe" -OutFile "dns-switcher.exe"
```

---

### Option 2: Install with Go

```bash
go install github.com/Har2yQn78/dns-switcher@latest
```

### Option 3: Build from Source

```bash
# Clone the repository
git clone https://github.com/Har2yQn78/dns-switcher.git
cd dns-switcher

# Build
go build -ldflags="-s -w" -o dns-switcher

# Install
sudo cp dns-switcher /usr/local/bin/
sudo chmod +x /usr/local/bin/dns-switcher
```

## Usage

**Linux/macOS:**

```bash
sudo dns-switcher
```

**Windows:**

```powershell
# Run as Administrator (Right-click PowerShell -> Run as Administrator)
.\dns-switcher.exe
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
```
