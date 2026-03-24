package internal

import (
	"crypto/tls"
	"errors"
	"net"

	"github.com/penndev/gopkg/socks5"
	"github.com/penndev/gopkg/util"
	"github.com/penndev/socks5/core/internal/dialer"
)

type HandleConnect func(net.Conn, net.Addr) error

// 本地请求，不用远程
func Local() HandleConnect {
	return func(conn net.Conn, addr net.Addr) error {
		var dial *net.Dialer
		switch addr.Network() {
		case "tcp":
			dial = dialer.TCPDialer
		case "udp":
			dial = dialer.UDPDialer
		}
		remote, err := dial.Dial(addr.Network(), addr.String())
		if err != nil {
			return err
		}
		util.Pipe(conn, remote)
		return nil
	}
}

// socks5标准请求
func Socks5(host, user, pass string) HandleConnect {
	return func(conn net.Conn, addr net.Addr) error {
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
		remote, err := socks.Dial(addr.Network(), addr.String())
		if err != nil {
			return err
		}
		util.Pipe(conn, remote)
		return nil
	}
}

// socks5 tls
func Socks5OverTLS(host, user, pass string, conf *tls.Config) HandleConnect {
	return func(conn net.Conn, addr net.Addr) error {
		s5tls, err := tls.Dial("tcp", host, conf)
		if err != nil {
			return errors.New("invalid rtype")
		}
		socks := &socks5.Client{
			Username: user,
			Password: pass,
			Conn:     s5tls,
		}
		err = socks.Negotiation()
		if err != nil {
			return err
		}
		// 发起真实的请求
		// 发起真实的请求
		remote, err := socks.Dial(addr.Network(), addr.String())
		if err != nil {
			return err
		}
		util.Pipe(conn, remote)
		return nil
	}
}
