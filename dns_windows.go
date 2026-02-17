//go:build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func powershell(args ...string) *exec.Cmd {
	pwsh := "C:\\Program Files\\PowerShell\\7\\pwsh.exe"
	if _, err := os.Stat(pwsh); err == nil {
		return exec.Command(pwsh, args...)
	}

	return exec.Command("C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe", args...)
}

func IsAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	return true
}

func GetCurrentDNS() ([]string, error) {
	adapter, err := getActiveNetworkAdapter()
	if err != nil {
		return nil, err
	}

	cmd := powershell("-Command",
		fmt.Sprintf("(Get-DnsClientServerAddress -InterfaceAlias '%s' -AddressFamily IPv4).ServerAddresses", adapter))
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get DNS servers: %w", err)
	}

	result := strings.TrimSpace(string(output))
	if result == "" {
		return []string{}, nil
	}

	lines := strings.Split(result, "\r\n")
	var dnsServers []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			dnsServers = append(dnsServers, line)
		}
	}

	return dnsServers, nil
}

func getActiveNetworkAdapter() (string, error) {
	cmd := powershell("-Command",
		"Get-NetAdapter | Where-Object {$_.Status -eq 'Up'} | Select-Object -First 1 -ExpandProperty Name")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get network adapter: %w", err)
	}

	adapter := strings.TrimSpace(string(output))
	if adapter == "" {
		return "", fmt.Errorf("no active network adapter found")
	}

	return adapter, nil
}

func BackupResolvConf() (string, error) {
	return "backup not needed on Windows", nil
}

func UpdateResolvConf(provider DNSProvider) error {
	adapter, err := getActiveNetworkAdapter()
	if err != nil {
		return err
	}

	var cmd *exec.Cmd

	if provider.Name == "Reset to Default" {
		cmd = powershell("-Command",
			fmt.Sprintf("Set-DnsClientServerAddress -InterfaceAlias '%s' -ResetServerAddresses", adapter))
	} else {
		servers := strings.Join(provider.Servers, ",")
		cmd = powershell("-Command",
			fmt.Sprintf("Set-DnsClientServerAddress -InterfaceAlias '%s' -ServerAddresses %s", adapter, servers))
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update DNS: %w (make sure you are running as Administrator)", err)
	}

	flushCmd := exec.Command("C:\\Windows\\System32\\ipconfig.exe", "/flushdns")
	_ = flushCmd.Run()

	return nil
}

func RestartSystemdResolved() error {
	return nil
}
