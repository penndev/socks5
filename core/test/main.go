package main

import (
	"fmt"
	"log"
	"net"
	"net/netip"
)

func main() {
	p := netip.MustParsePrefix("192.168.1.123/24")

	ip := p.Addr()

	if ip.Is4() {
		fmt.Println("IPv4")
	} else {
		fmt.Println("IPv6")
	}

	fmt.Println("IP:", ip)
	fmt.Println("Mask bits:", net.IP(net.CIDRMask(p.Bits(), 32)).String())

	// prefix.Addr().String(),
	// prefix.Masked().Addr().String(),

	conn, _ := net.Dial("udp", "8.8.8.8:80")
	lip := conn.LocalAddr().(*net.UDPAddr).IP.String()
	conn.Close()
	log.Println(lip)

}
