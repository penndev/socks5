package proxy

import (
	"net/netip"

	"golang.org/x/sys/windows"
	"golang.zx2c4.com/wireguard/tun"
)

const TUN_NAME = "prise-tun"
const TUN_MTU = 0
const TUN_OFFSET = 0

var TUN_IP netip.Prefix
var Routes []netip.Prefix

// 自定义网卡GUID 方便wintun复用
func init() {
	TUN_IP = netip.MustParsePrefix("172.19.0.1/32")
	Routes = []netip.Prefix{netip.MustParsePrefix("0.0.0.0/0")}
	// 设置tun设备名称标识和guid
	tun.WintunTunnelType = TUN_NAME
	tun.WintunStaticRequestedGUID = &windows.GUID{
		Data1: 0x8ceeab57,
		Data2: 0x7cb2,
		Data3: 0x469f,
		Data4: [8]byte{0x91, 0x3b, 0xea, 0xeb, 0x22, 0xe2, 0x28, 0x24},
	}
}
