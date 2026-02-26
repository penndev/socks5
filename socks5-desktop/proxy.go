package main

import (
	"crypto/tls"
	"net"

	"github.com/penndev/gopkg/socks5"
)

type Proxy struct {
	// 远程socks5服务器信息
	rhost       string
	ruser       string
	rpass       string
	rtype       string         //判断是 Socks5 还是 Socks5OverTls
	localServer *socks5.Server //本地socks5服务
}

func (p *Proxy) HandleConnect(conn net.Conn, req socks5.Requests, replies func(status socks5.REP) error) error {
	addr := req.Addr()

	var (
		remote net.Conn
		err    error
	)

	// 根据 rtype 判断是否走 TLS
	if p.rtype == "Socks5OverTls" {
		tlsConf := &tls.Config{
			InsecureSkipVerify: true,
		}
		remote, err = tls.Dial("tcp", addr, tlsConf)
	} else {
		remote, err = net.Dial("tcp", addr)
	}

	if err != nil {
		replies(socks5.REP_CONNECTION_REFUSED)
		return err
	}

	replies(socks5.REP_SUCCEEDED)
	defer remote.Close()
	socks5.Pipe(conn, remote)
	return nil
}

func (p *Proxy) Start(host, user, pass string) error {

	p.localServer = &socks5.Server{
		Addr:          host,
		Username:      user,
		Password:      pass,
		HandleConnect: p.HandleConnect,
	}
	if p.localServer.Username != "" {
		p.localServer.METHOD = socks5.METHOD_USERNAME_PASSWORD
	}
	go func() {
		if err := p.localServer.TCPListen(); err != nil {

		}
	}()
	return nil
}

func (p *Proxy) Stop() error {
	p.localServer.Close()
	return nil
}

func (p *Proxy) SetRemote(host, user, pass, rtype string) error {
	p.rhost = host
	p.ruser = user
	p.rpass = pass
	p.rtype = rtype
	return nil
}
