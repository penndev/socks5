package main

import (
	"net"
	"net/url"

	"github.com/penndev/gopkg/socks5"
)

type Proxy struct {
	remote      *url.URL       // 远程连接信息
	localServer *socks5.Server //本地socks5服务
}

func (p *Proxy) setLocalConnect(c net.Conn, req socks5.Requests, rep socks5.HandleReply) error {
	addr := req.Addr()
	app.Event.Emit(appConst.LogTypeName_LOG, "incoming: "+addr)
	defer app.Event.Emit(appConst.LogTypeName_LOG, "proxy finished: "+addr)

	return nil
}

func (p *Proxy) SetLocal(host, user, pass string) error {
	app.Event.Emit(
		appConst.LogTypeName_STATUS,
		"localServer://"+user+":"+pass+"@"+host,
	)
	if p.localServer != nil {
		// 配置未变化，保持当前服务
		if p.localServer.Addr == host &&
			p.localServer.Username == user &&
			p.localServer.Password == pass {
			return nil
		}
		// 配置变化，先关闭旧的监听
		p.localServer.Close()
		p.localServer = nil
	}

	p.localServer = &socks5.Server{
		Addr:          host,
		Username:      user,
		Password:      pass,
		HandleConnect: p.setLocalConnect,
	}
	go func() {
		if err := p.localServer.TCPListen(); err != nil {
			app.Event.Emit(
				appConst.LogTypeName_STATUS,
				"local server start failed: "+err.Error(),
			)
		}
	}()
	return nil
}

func (p *Proxy) SetRemote(remote string) error {
	ru, err := url.Parse(remote)
	if err != nil {
		app.Event.Emit(appConst.LogTypeName_STATUS, err.Error())
	}

	app.Event.Emit(
		appConst.LogTypeName_STATUS,
		"SetRemote-> "+ru.Scheme+"://"+ru.User.String()+"@"+ru.Host,
	)
	p.remote = ru
	return nil
}

func (p *Proxy) SetMode(mode string) error {
	return nil
}
