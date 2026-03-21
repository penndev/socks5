package tun

import "net/netip"

type Options struct {
	Name string
	MTU  int
	IP   netip.Prefix
}
