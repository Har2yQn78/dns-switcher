package main

import (
	"context"
	"fmt"
	"net"
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
