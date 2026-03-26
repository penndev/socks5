package transport

import (
	"crypto/tls"
	"net"

	"github.com/penndev/gopkg/socks5"
	"github.com/penndev/gopkg/util"
	"github.com/penndev/socks5/core/transport/dialer"
)

type HandleConnect func(net.Conn, string, string) error

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

// socks5标准请求
func Socks5(host, user, pass string) HandleConnect {
	return func(conn net.Conn, network, address string) error {
		// 控制直接走物理网卡防止环路。如果是本地可能失败，如果需要处理则判断host是否走 ip.IsLoopback
		dialTcp, err := dialer.TCPDialer.Dial("tcp", host)
		if err != nil {
			return err
		}
		socks := &socks5.Client{
			Username: user,
			Password: pass,
			Conn:     dialTcp,
		}
		// socks5 握手
		err = socks.Negotiation()
		if err != nil {
			return err
		}
		// 发起真实的请求
		remote, err := socks.Dial(network, address)
		if err != nil {
			return err
		}
		util.Pipe(conn, remote)
		return nil
	}
}

// socks5 tls
func Socks5OverTLS(host, user, pass string, conf *tls.Config) HandleConnect {
	return func(conn net.Conn, network, address string) error {
		// s5tls, err := tls.Dial("tcp", host, conf)
		// if err != nil {
		// 	return err
		// }
		tcpDial, err := dialer.TCPDialer.Dial("tcp", host)
		if err != nil {
			return err
		}
		conf.InsecureSkipVerify = true
		tlsConn := tls.Client(tcpDial, conf)
		if err = tlsConn.Handshake(); err != nil {
			return err
		}
		socks := &socks5.Client{
			Username: user,
			Password: pass,
			Conn:     tlsConn,
		}
		err = socks.Negotiation()
		if err != nil {
			return err
		}
		remote, err := socks.Dial(network, address)
		if err != nil {
			return err
		}
		util.Pipe(conn, remote)
		return nil
	}
}
