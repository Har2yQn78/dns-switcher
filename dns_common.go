package main

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"
	"time"
)

func ValidateDNS(servers []string) (bool, error) {
	testDomain := "google.com"

	for _, server := range servers {
		resolver := &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: 3 * time.Second,
				}
				return d.Dial(network, net.JoinHostPort(server, "53"))
			},
		}

		ctx := context.Background()
		_, err := resolver.LookupHost(ctx, testDomain)

		if err != nil {
			return false, fmt.Errorf("DNS server %s is not responding", server)
		}
	}

	return true, nil
}

func TestDNSLatency(server string) int {
	testDomain := "google.com"

	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 2 * time.Second,
			}
			return d.Dial(network, net.JoinHostPort(server, "53"))
		},
	}

	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := resolver.LookupHost(ctx, testDomain)

	elapsed := time.Since(start).Milliseconds()

	if err != nil {
		return -1
	}

	return int(elapsed)
}

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
