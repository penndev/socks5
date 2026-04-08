package dialer

import "net"

// 本地网卡IP 解决环路
var LocalIP net.IP

var TCPDialer *net.Dialer

var UDPDialer *net.Dialer

func init() {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()

	LocalIP = conn.LocalAddr().(*net.UDPAddr).IP

	TCPDialer = &net.Dialer{
		LocalAddr: &net.TCPAddr{
			IP: LocalIP,
		},
	}
	UDPDialer = &net.Dialer{
		LocalAddr: &net.UDPAddr{
			IP: LocalIP,
		},
	}
}
