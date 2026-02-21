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
	{Name: "AdGuard", Servers: []string{"94.140.14.14", "94.140.15.15"}, Latency: -1},
	{Name: "Quad9", Servers: []string{"9.9.9.9", "149.112.112.112"}, Latency: -1},
	{Name: "OpenDNS", Servers: []string{"208.67.222.222", "208.67.220.220"}, Latency: -1},
	{Name: "Level3", Servers: []string{"4.2.2.1", "4.2.2.2"}, Latency: -1},
	{Name: "Verisign", Servers: []string{"64.6.64.6", "64.6.65.6"}, Latency: -1},
	{Name: "UltraDNS", Servers: []string{"156.154.70.1", "156.154.71.1"}, Latency: -1},
	{Name: "DNS.WATCH", Servers: []string{"84.200.69.80", "84.200.70.40"}, Latency: -1},
	{Name: "Comodo", Servers: []string{"8.26.56.26", "8.20.247.20"}, Latency: -1},
	{Name: "CleanBrowsing", Servers: []string{"185.228.168.9", "185.228.169.9"}, Latency: -1},
	{Name: "Neustar", Servers: []string{"156.154.70.2", "156.154.71.2"}, Latency: -1},
	{Name: "Yandex.DNS", Servers: []string{"77.88.8.8", "77.88.8.1"}, Latency: -1},
	{Name: "Freenom World", Servers: []string{"80.80.80.80", "80.80.81.81"}, Latency: -1},
	{Name: "Reset to Default", Servers: []string{"127.0.0.53"}, Latency: -1},
	{Name: "Add Custom DNS", Servers: []string{}, Latency: -1},
}
