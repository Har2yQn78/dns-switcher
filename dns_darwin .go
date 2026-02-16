// +build darwin

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func IsAdmin() bool {
	return os.Geteuid() == 0
}

func GetCurrentDNS() ([]string, error) {
	// Get the active network service
	service, err := getActiveNetworkService()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("networksetup", "-getdnsservers", service)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get DNS servers: %w", err)
	}

	result := strings.TrimSpace(string(output))
	if result == "There aren't any DNS Servers set on "+service+"." {
		return []string{}, nil
	}

	dnsServers := strings.Split(result, "\n")
	return dnsServers, nil
}

func getActiveNetworkService() (string, error) {
	// Try to get Wi-Fi first
	cmd := exec.Command("networksetup", "-listallnetworkservices")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to list network services: %w", err)
	}

	services := strings.Split(string(output), "\n")
	for _, service := range services {
		service = strings.TrimSpace(service)
		if service == "" || strings.HasPrefix(service, "*") {
			continue
		}

		if strings.Contains(service, "Wi-Fi") || strings.Contains(service, "Ethernet") {
			return service, nil
		}
	}

	for _, service := range services {
		service = strings.TrimSpace(service)
		if service != "" && !strings.HasPrefix(service, "*") && !strings.HasPrefix(service, "An asterisk") {
			return service, nil
		}
	}

	return "", fmt.Errorf("no active network service found")
}

func BackupResolvConf() (string, error) {
	return "backup not needed on macOS", nil
}

func UpdateResolvConf(provider DNSProvider) error {
	service, err := getActiveNetworkService()
	if err != nil {
		return err
	}

	args := []string{"-setdnsservers", service}
	if provider.Name == "Reset to Default" {
		args = append(args, "empty")
	} else {
		args = append(args, provider.Servers...)
	}

	cmd := exec.Command("networksetup", args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update DNS: %w", err)
	}

	return nil
}

func RestartSystemdResolved() error {
	return nil
}
