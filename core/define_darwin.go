package main

import "net/netip"

// netsh interface ipv4 set address name="PrismTUN" source=static addr=172.19.0.1 mask=255.255.255.255
const TUN_NAME = "utun"
const TUN_MTU = 1500
const TUN_OFFSET = 4

var TUN_IP netip.Prefix
var Routes []netip.Prefix

// 自定义网卡GUID 方便wintun复用
func init() {
	TUN_IP = netip.MustParsePrefix("172.19.0.1/32")
	Routes = []netip.Prefix{
		netip.MustParsePrefix("1.0.0.0/8"),
		netip.MustParsePrefix("2.0.0.0/7"),
		netip.MustParsePrefix("4.0.0.0/6"),
		netip.MustParsePrefix("8.0.0.0/5"),
		netip.MustParsePrefix("16.0.0.0/4"),
		netip.MustParsePrefix("32.0.0.0/3"),
		netip.MustParsePrefix("64.0.0.0/2"),
		netip.MustParsePrefix("128.0.0.0/1"),
		netip.MustParsePrefix("198.18.0.0/15"),
	}
}
