package main

import (
	"crypto/tls"
	"errors"
	"net"

	"github.com/penndev/gopkg/socks5"
	"github.com/penndev/socks5/core/stack"
	"github.com/penndev/socks5/core/tun"
)

type Proxy struct {
	// 远程socks5服务器信息
	rhost       string
	ruser       string
	rpass       string
	rtype       string         //判断是 Socks5 还是 Socks5OverTLS
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
	switch p.rtype {
	case "Socks5OverTLS":
		tlsConf := &tls.Config{
			InsecureSkipVerify: true,
		}
		server, err = tls.Dial("tcp", p.rhost, tlsConf)
	case "Socks5":
		server, err = net.Dial("tcp", p.rhost)
	default:
		app.Event.Emit("logProxyList", "invalid rtype: "+p.rtype)
		return errors.New("invalid rtype")
	}
	if err != nil {
		app.Event.Emit("logProxyList", "connect upstream failed: "+err.Error())

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
	// 已经启动过则检查配置是否有变化：
	// - 配置未变化：什么也不做
	// - 配置有变化：关闭旧服务后按新配置重启
	if p.localServer != nil {
		if p.localServer.Addr == host &&
			p.localServer.Username == user &&
			p.localServer.Password == pass {
			// 配置未变化，保持当前服务
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
		HandleConnect: p.HandleConnect,
	}
	if p.localServer.Username != "" {
		p.localServer.METHOD = socks5.METHOD_USERNAME_PASSWORD
	} else {
		p.localServer.METHOD = socks5.METHOD_NO_AUTH
	}
	go func() {
		if err := p.localServer.TCPListen(); err != nil {
			app.Event.Emit("logServerStatus", "local server start failed: "+err.Error()+"\n")
		}
	}()
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

func (p *Proxy) SetMode(mode string) error {
	app.Event.Emit("logServerStatus", "set mode: "+mode+"\n")
	switch mode {
	case "manual":
		p.systemStop()
	case "system":
		err := p.systemStart()
		if err != nil {
			app.Event.Emit("logServerStatus", "system start: "+err.Error()+"\n")
		}
		return nil
	case "tun":
		// todo: tun mode
	default:
		return errors.New("invalid mode")
	}
	return nil
}

func (p *Proxy) setModeTun() error {
	dev, err := tun.CreateTUN("dev_socks5", 0)
	if err != nil {
		return err
	}

	stack.New(stack.Option{
		EndPoint: dev,
		HandleTCP: func(ftr *stack.ForwarderTCPRequest) {
			ftr.Conn

			// log.Printf("tcp %s <-> %s", ftr.LocalAddr, ftr.RemoteAddr)
			// tunnel.Tunnel(tunnel.Option{
			// 	Src:    ftr.Conn,
			// 	Dst:    remoteConn,
			// 	BufLen: 32 * 1024,
			// })
		},
		// HandlerUDP: func(fur *stack.ForwarderUDPRequest) {
		// 	defer fur.Conn.Close()
		// 	s5, err := socks5.NewClient(host, user, pass)
		// 	if err != nil {
		// 		log.Println("socks5 connection err:", err)
		// 		return
		// 	}
		// 	defer s5.Close()

		// 	remoteConn, err := s5.Dial("udp", fur.RemoteAddr)
		// 	if err != nil {
		// 		log.Println("socks5 udp remote err:", err)
		// 		return
		// 	}
		// 	log.Printf("udp %s <-> %s", fur.LocalAddr, fur.RemoteAddr)
		// 	tunnel.Tunnel(tunnel.Option{
		// 		Src:     fur.Conn,
		// 		Dst:     remoteConn,
		// 		BufLen:  1024,
		// 		Timeout: 30 * time.Second,
		// 	})
		// },
	})

}
