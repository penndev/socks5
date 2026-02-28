package main

import (
	"crypto/tls"
	"log"
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
	app.Event.Emit("logProxyList", "incoming: "+addr)
	defer app.Event.Emit("logProxyList", "proxy finished: "+addr)
	var (
		server net.Conn
		err    error
	)
	// 根据 rtype 判断是否走 TLS
	if p.rtype == "Socks5OverTls" {
		tlsConf := &tls.Config{
			InsecureSkipVerify: true,
		}
		server, err = tls.Dial("tcp", p.rhost, tlsConf)
	} else {
		server, err = net.Dial("tcp", p.rhost)
	}
	if err != nil {
		app.Event.Emit("logProxyList", "connect upstream failed: "+err.Error())
	}
	if err != nil {
		replies(socks5.REP_CONNECTION_REFUSED)
		return err
	}
	socks := &socks5.Client{
		Username: p.ruser,
		Password: p.rpass,
		Conn:     server,
	}
	// socks5 握手
	err = socks.Negotiation()
	if err != nil {
		app.Event.Emit("logProxyList", "socks5 negotiation failed: "+err.Error())
		replies(socks5.REP_CONNECTION_REFUSED)
		return err
	}
	// 发起真实的请求
	remote, err := socks.Dial("tcp", addr)
	if err != nil {
		app.Event.Emit("logProxyList", "dial target failed: "+err.Error())
		replies(socks5.REP_CONNECTION_REFUSED)
		return err
	}
	replies(socks5.REP_SUCCEEDED)
	defer remote.Close()
	socks5.Pipe(conn, remote)
	return nil
}

func (p *Proxy) Start(host, user, pass string) error {
	app.Event.Emit("logServerStatus", "local: socks5://"+user+":"+pass+"@"+host+"\n")
	p.localServer = &socks5.Server{
		Addr:          host,
		Username:      user,
		Password:      pass,
		HandleConnect: p.HandleConnect,
	}
	if p.localServer.Username != "" {
		p.localServer.METHOD = socks5.METHOD_USERNAME_PASSWORD
	} else {
		p.localServer.METHOD = socks5.METHOD_NO_AUTH
	}
	go func() {
		log.Println("socks5 启动成功")
		if err := p.localServer.TCPListen(); err != nil {
			app.Event.Emit("logServerStatus", "local server start failed: "+err.Error()+"\n")
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

	app.Event.Emit("logServerStatus", "remote: "+rtype+"://"+user+":"+pass+"@"+host+"\n")
	return nil
}
