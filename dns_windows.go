//go:build windows

package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetCurrentDNS() ([]string, error) {
	adapter, err := getActiveNetworkAdapter()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("powershell", "-Command",
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
	cmd := exec.Command("powershell", "-Command",
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
		cmd = exec.Command("powershell", "-Command",
			fmt.Sprintf("Set-DnsClientServerAddress -InterfaceAlias '%s' -ResetServerAddresses", adapter))
	} else {
		servers := strings.Join(provider.Servers, ",")
		cmd = exec.Command("powershell", "-Command",
			fmt.Sprintf("Set-DnsClientServerAddress -InterfaceAlias '%s' -ServerAddresses %s", adapter, servers))
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update DNS: %w (you may need to run as Administrator)", err)
	}

	flushCmd := exec.Command("ipconfig", "/flushdns")
	_ = flushCmd.Run()

	return nil
}

func RestartSystemdResolved() error {
	return nil
}
