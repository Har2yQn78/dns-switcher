package main

type DNSProvider struct {
	Name    string
	Servers []string
}

var providers = []DNSProvider{
	{Name: "Shecan", Servers: []string{"178.22.122.100", "185.51.200.2"}},
	{Name: "Radar", Servers: []string{"10.202.10.10", "10.202.10.11"}},
	{Name: "Electro", Servers: []string{"78.157.42.100", "78.157.42.101"}},
	{Name: "Begzar", Servers: []string{"185.55.226.26", "185.55.226.25"}},
	{Name: "UltraDNS", Servers: []string{"64.6.64.6", "64.6.65.6"}},
	{Name: "DNS Pro", Servers: []string{"87.107.110.109", "87.107.110.110"}},
	{Name: "DynX", Servers: []string{"10.70.95.150", "10.70.95.162"}},
	{Name: "403", Servers: []string{"10.202.10.202", "10.202.10.102"}},
	{Name: "Google", Servers: []string{"8.8.8.8", "8.8.4.4"}},
	{Name: "Cloudflare", Servers: []string{"1.1.1.1", "1.0.0.1"}},
	{Name: "Reset to Default", Servers: []string{"127.0.0.53"}},
}
