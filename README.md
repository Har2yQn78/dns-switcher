# DNS Changer

A fast, simple CLI tool to change DNS servers on Linux systems with a TUI interface.

![DNS Changer](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

## Features

- Fast and lightweight
- Automatic backup of DNS configuration (Linux)
- DNS validation after changes
- Automatic systemd-resolved restart (Linux)
- Single binary, no dependencies
- Linux support
- macOS support (Intel & Apple Silicon)
- Windows support (64-bit, 32-bit, ARM64)

## Supported DNS Providers

- Shecan
- Radar
- Electro
- Begzar
- UltraDNS
- DNS Pro
- DynX
- 403
- Google DNS
- Cloudflare DNS
- Reset to Default

## Installation

### Option 1: Install from Release (Recommended)

Download the latest release from the [releases page](https://github.com/Har2yQn78/dns-switcher/releases):

**Linux:**

```bash
# Download the binary
wget https://github.com/Har2yQn78/dns-switcher/releases/download/v1.0.0/dns-changer-linux-amd64

# Install it
sudo mv dns-changer-linux-amd64 /usr/local/bin/dns-changer
sudo chmod +x /usr/local/bin/dns-changer
```

**macOS (Intel):**

```bash
# Download the binary
curl -L https://github.com/Har2yQn78/dns-switcher/releases/latest/download/dns-changer-macos-amd64 -o dns-changer

# Install it
sudo mv dns-changer /usr/local/bin/
sudo chmod +x /usr/local/bin/dns-changer
```

**macOS (Apple Silicon):**

```bash
# Download the binary
curl -L https://github.com/Har2yQn78/dns-switcher/releases/latest/download/dns-changer-macos-arm64 -o dns-changer

# Install it
sudo mv dns-changer /usr/local/bin/
sudo chmod +x /usr/local/bin/dns-changer
```

**Windows (64-bit):**

```powershell
# Download using PowerShell
Invoke-WebRequest -Uri "https://github.com/Har2yQn78/dns-switcher/releases/latest/download/dns-changer-windows-amd64.exe" -OutFile "dns-changer.exe"

# Move to a directory in your PATH (optional)
# Or run directly from current directory
```

**Windows (32-bit):**

```powershell
# Download using PowerShell
Invoke-WebRequest -Uri "https://github.com/Har2yQn78/dns-switcher/releases/latest/download/dns-changer-windows-386.exe" -OutFile "dns-changer.exe"
```

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

**Note:** Administrator/root privileges are required to modify DNS settings.

### Navigation

- `↑/↓` or `j/k` - Navigate through providers
- `Enter` - Select a provider
- `q` or `Ctrl+C` - Quit

## Requirements

- **Linux** (tested on Ubuntu/Pop!\_OS), **macOS** (10.13+), or **Windows** (10/11)
- Administrator/root/sudo access
- Go 1.21+ (for building from source)

## How It Works

**Linux:**

1. Displays current DNS servers
2. Shows a list of available DNS providers
3. Backs up current `/etc/resolv.conf`
4. Updates DNS configuration
5. Restarts `systemd-resolved` if active
6. Validates DNS servers are responding

**macOS:**

1. Displays current DNS servers
2. Shows a list of available DNS providers
3. Uses `networksetup` to change DNS for active network service
4. Validates DNS servers are responding

**Windows:**

1. Displays current DNS servers
2. Shows a list of available DNS providers
3. Uses PowerShell `Set-DnsClientServerAddress` to update DNS
4. Flushes DNS cache
5. Validates DNS servers are responding

## Backup

**Linux only:** Every time you change DNS, a backup is created at:

```
/etc/resolv.conf.bak.YYYYMMDD_HHMMSS
```

To restore from backup:

```bash
sudo cp /etc/resolv.conf.bak.YYYYMMDD_HHMMSS /etc/resolv.conf
```

**macOS:** DNS changes can be reverted by selecting "Reset to Default" or manually via System Preferences → Network.

**Windows:** DNS changes can be reverted by selecting "Reset to Default" or manually via Settings → Network & Internet → Change adapter options → Properties → IPv4 Properties.

## Configuration

All DNS providers are defined in `providers.go`. To add your own:

```go
{Name: "Custom DNS", Servers: []string{"1.2.3.4", "5.6.7.8"}},
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Author

Created by [Harry](https://github.com/Har2yQn78)

## Acknowledgments

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- Styled with [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
