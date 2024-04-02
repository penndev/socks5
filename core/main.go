package main

import "flag"

type Config struct {
	Proxy   string
	TunName string
	TunIP   string
	TunMtu  int
}

var (
	versionFlag = false
	version     = "0.0.1"
	config      = Config{}
)

func init() {
	flag.StringVar(&config.Proxy, "proxy", "", "Set remote socks5 service example socks5://user:pass@192.168.0.1:1080")
	flag.StringVar(&config.TunIP, "ip", "10.10.10.10", "Set device static ip")
	flag.IntVar(&config.TunMtu, "mtu", 0, "Set tun device mtu")
	flag.StringVar(&config.TunName, "name", "socks5", "Set device name")
	flag.BoolVar(&versionFlag, "version", false, "Show version and then quit")
	flag.Parse()
}
