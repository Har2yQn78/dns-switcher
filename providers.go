package main

type DNSProvider struct {
	Name    string
	Servers []string
	Latency int
}

var providers = []DNSProvider{
	{Name: "Shecan", Servers: []string{"178.22.122.100", "185.51.200.2"}, Latency: -1},
	{Name: "Radar", Servers: []string{"10.202.10.10", "10.202.10.11"}, Latency: -1},
	{Name: "Electro", Servers: []string{"78.157.42.100", "78.157.42.101"}, Latency: -1},
	{Name: "Begzar", Servers: []string{"185.55.226.26", "185.55.226.25"}, Latency: -1},
	{Name: "DNS Pro", Servers: []string{"87.107.110.109", "87.107.110.110"}, Latency: -1},
	{Name: "DynX", Servers: []string{"10.70.95.150", "10.70.95.162"}, Latency: -1},
	{Name: "403", Servers: []string{"10.202.10.202", "10.202.10.102"}, Latency: -1},
	{Name: "Google", Servers: []string{"8.8.8.8", "8.8.4.4"}, Latency: -1},
	{Name: "Cloudflare", Servers: []string{"1.1.1.1", "1.0.0.1"}, Latency: -1},
	{Name: "Reset to Default", Servers: []string{"127.0.0.53"}, Latency: -1},
	{Name: "Add Custom DNS", Servers: []string{}, Latency: -1},
}
