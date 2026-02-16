package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const resolvConfPath = "/etc/resolv.conf"

func GetCurrentDNS() ([]string, error) {
	file, err := os.Open(resolvConfPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", resolvConfPath, err)
	}
	defer file.Close()

	var dnsServers []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "nameserver") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				dnsServers = append(dnsServers, fields[1])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return dnsServers, nil
}

func BackupResolvConf() (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.bak.%s", resolvConfPath, timestamp)

	input, err := os.ReadFile(resolvConfPath)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", resolvConfPath, err)
	}

	err = os.WriteFile(backupPath, input, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create backup: %w", err)
	}

	return backupPath, nil
}

func UpdateResolvConf(provider DNSProvider) error {
	backupPath, err := BackupResolvConf()
	if err != nil {
		return fmt.Errorf("backup failed: %w", err)
	}

	var content strings.Builder
	content.WriteString(fmt.Sprintf("# Updated on %s\n", time.Now().Format("2006-01-02 15:04:05")))

	if provider.Name != "Reset to Default" {
		content.WriteString(fmt.Sprintf("# Provider: %s\n", provider.Name))
	}

	for _, dns := range provider.Servers {
		content.WriteString(fmt.Sprintf("nameserver %s\n", dns))
	}

	if provider.Name != "Reset to Default" {
		content.WriteString("options edns0 trust-ad\n")
	}

	err = os.WriteFile(resolvConfPath, []byte(content.String()), 0644)
	if err != nil {
		return fmt.Errorf("failed to write %s: %w", resolvConfPath, err)
	}

	fmt.Printf("✓ Backup created: %s\n", backupPath)
	return nil
}

func RestartSystemdResolved() error {
	cmd := exec.Command("systemctl", "is-active", "systemd-resolved")
	err := cmd.Run()

	if err == nil {
		cmd = exec.Command("systemctl", "restart", "systemd-resolved")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to restart systemd-resolved: %w", err)
		}
		fmt.Println("✓ systemd-resolved restarted")
	}

	return nil
}
