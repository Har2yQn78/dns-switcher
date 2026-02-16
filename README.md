# DNS Changer

A fast, simple CLI tool to change DNS servers on Linux systems with a beautiful TUI interface.

![DNS Changer](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

## Features

- Fast and lightweight
- Automatic backup of DNS configuration
- DNS validation after changes
- Automatic systemd-resolved restart
- Single binary, no dependencies

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

## Installation

### Option 1: Install from Release (Recommended)

Download the latest release from the [releases page](https://github.com/Har2yQn78/dns-switcher/releases):

```bash
# Download the binary
wget https://github.com/Har2yQn78/dns-switcher/releases/download/v1.0.0/dns-changer

# Install it
sudo chmod +x /usr/local/bin/dns-changer
```

### Option 2: Install with Go

```bash
go install github.com/Har2yQn78/dns-switcher
```

### Option 3: Build from Source

```bash
# Clone the repository
git clone https://github.com/Har2yQn78/dns-switcher.git
cd dns-changer

# Build
go build -ldflags="-s -w" -o dns-changer

# Install
sudo cp dns-changer /usr/local/bin/
sudo chmod +x /usr/local/bin/dns-changer
```

## Usage

Simply run:

```bash
sudo dns-changer
```

**Note:** Root privileges are required to modify DNS settings.

### Navigation

- `↑/↓` or `j/k` - Navigate through providers
- `Enter` - Select a provider
- `q` or `Ctrl+C` - Quit

## Requirements

- Linux
- Root/sudo access
- Go 1.21+ (for building from source)

## How It Works

1. Displays current DNS servers
2. Shows a list of available DNS providers
3. Backs up current `/etc/resolv.conf`
4. Updates DNS configuration
5. Restarts `systemd-resolved` if active
6. Validates DNS servers are responding

## Backup

Every time you change DNS, a backup is created at:

```
/etc/resolv.conf.bak.YYYYMMDD_HHMMSS
```

To restore from backup:

```bash
sudo cp /etc/resolv.conf.bak.YYYYMMDD_HHMMSS /etc/resolv.conf
```

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
