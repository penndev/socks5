package tun

import "net/netip"

type Options struct {
	Name   string
	MTU    int
	Offset int
	IP     netip.Prefix
}
