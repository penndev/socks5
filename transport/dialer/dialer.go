package dialer

import (
	"net"
)

var TCPDialer *net.Dialer = &net.Dialer{}

var UDPDialer *net.Dialer = &net.Dialer{}
