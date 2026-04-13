package transport

import (
	"net"

	"github.com/penndev/gopkg/util"
	"github.com/penndev/prism/transport/dialer"
)

// 本地请求，不用远程
func Local() HandleConnect {
	return func(conn net.Conn, network, address string) error {
		var dial *net.Dialer
		switch network {
		case "tcp":
			dial = dialer.TCPDialer
		case "udp":
			dial = dialer.UDPDialer
		}
		remote, err := dial.Dial(network, address)
		if err != nil {
			return err
		}
		util.Pipe(conn, remote)
		return nil
	}
}
