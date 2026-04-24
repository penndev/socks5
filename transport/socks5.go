package transport

import (
	"crypto/tls"
	"net"

	"github.com/penndev/gopkg/socks5"
	"github.com/penndev/gopkg/util"
	"github.com/penndev/prism/transport/dialer"
)

// socks5标准请求
func Socks5(host, user, pass string) HandleConnect {
	return func(conn net.Conn, network, address string) error {
		var dialTcp net.Conn
		var err error
		if isLoopback(host) {
			dialTcp, err = net.Dial("tcp", host)
		} else {
			dialTcp, err = dialer.TCPDialer.Dial("tcp", host)
		}
		if err != nil {
			return err
		}
		socks := &socks5.Client{
			Username: user,
			Password: pass,
			Conn:     dialTcp,
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

// socks5 tls
func Socks5OverTLS(host, user, pass string, conf *tls.Config) HandleConnect {
	return func(conn net.Conn, network, address string) error {
		var dialTcp net.Conn
		var err error
		if isLoopback(host) {
			dialTcp, err = net.Dial("tcp", host)
		} else {
			dialTcp, err = dialer.TCPDialer.Dial("tcp", host)
		}
		if err != nil {
			return err
		}
		dialTls := tls.Client(dialTcp, conf)
		if err = dialTls.Handshake(); err != nil {
			return err
		}
		socks := &socks5.Client{
			Username: user,
			Password: pass,
			Conn:     dialTls,
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
