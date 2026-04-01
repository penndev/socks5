package proxy

import (
	"log"
	"net"
	"net/netip"
	"net/url"
	"socks5-desktop/internal"

	"github.com/penndev/gopkg/socks5"
	"github.com/penndev/socks5/core/route"
	"github.com/penndev/socks5/core/stack"
	"github.com/penndev/socks5/core/tun"
)

type Proxy struct {
	remote      *url.URL       // 远程连接信息
	localServer *socks5.Server //本地socks5服务

	// tun用
	dev *tun.Tun
}

func (p *Proxy) setLocalConnect(c net.Conn, req socks5.Requests, rep socks5.HandleReply) error {
	host := req.Addr()

	network := ""
	switch req.CMD {
	case socks5.CMD_CONNECT:
		network = "tcp"
	case socks5.CMD_UDP_ASSOCIATE:
		network = "udp"
	default:
		rep(socks5.REP_COMMAND_NOT_SUPPORTED)
	}

	var err error
	// handle := transport.Local()
	handle, err := HandleConnect(p.remote)
	if err != nil {
		internal.App.Event.Emit(internal.AppConfig.LogTypeName_LOG, host+":"+err.Error())
		return err
	}
	// internal.App.Event.Emit(internal.AppConfig.LogTypeName_LOG, "incoming: "+network+" "+host)
	rep(socks5.REP_SUCCEEDED)
	err = handle(c, network, host)
	if err != nil {
		internal.App.Event.Emit(internal.AppConfig.LogTypeName_LOG, host+":"+err.Error())
	}
	return err
}

func (p *Proxy) setModeTun() error {
	var err error
	p.dev, err = tun.New(tun.Options{
		Name:   TUN_NAME,
		MTU:    TUN_MTU,
		Offset: TUN_OFFSET,
	})
	if err != nil {
		return err
	}
	// defer dev.Close()
	handle, err := HandleConnect(p.remote)
	if err != nil {
		return err
	}

	stack.New(stack.Option{
		EndPoint: p.dev,
		HandleTCP: func(f *stack.ForwarderTCPRequest) {
			log.Printf(
				"%s -> %s",
				f.RemoteAddr.Network(),
				f.RemoteAddr.String(),
			)
			handle(f.Conn, f.RemoteAddr.Network(), f.RemoteAddr.String())
		},
		HandlerUDP: func(f *stack.ForwarderUDPRequest) {
			log.Printf(
				"%s -> %s",
				f.RemoteAddr.Network(),
				f.RemoteAddr.String(),
			)
			handle(f.Conn, f.RemoteAddr.Network(), f.RemoteAddr.String())
		},
	})
	route.Start(route.Options{
		DevName:      p.dev.Name(),
		DevIP:        netip.MustParsePrefix("172.19.0.1/32"),
		RouteAddress: []netip.Prefix{netip.MustParsePrefix("0.0.0.0/0")},
	})
	return nil
}

func (p *Proxy) SetStart(host, user, pass string) error {
	internal.App.Event.Emit(
		internal.AppConfig.LogTypeName_STATUS,
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
		Addr:     host,
		Username: user,
		Password: pass,
		// 本地连接处理函数 入口
		HandleConnect: p.setLocalConnect,
	}
	go func() {
		if err := p.localServer.TCPListen(); err != nil {
			internal.App.Event.Emit(
				internal.AppConfig.LogTypeName_STATUS,
				"p.localServer.TCPListen error: "+err.Error(),
			)
		}
	}()
	return nil
}

func (p *Proxy) SetRemote(remote string) error {
	ru, err := url.Parse(remote)
	if err != nil {
		internal.App.Event.Emit(internal.AppConfig.LogTypeName_STATUS, err.Error())
	}

	internal.App.Event.Emit(
		internal.AppConfig.LogTypeName_STATUS,
		"SetRemote-> "+ru.Scheme+"://"+ru.User.String()+"@"+ru.Host,
	)
	p.remote = ru
	return nil
}

func (p *Proxy) SetMode(mode string) error {
	internal.App.Event.Emit(
		internal.AppConfig.LogTypeName_STATUS,
		"SetMode-> "+mode,
	)
	switch mode {
	case "tun":
		return p.setModeTun()
	default:
		if p.dev != nil {
			p.dev.Close()
		}
	}
	return nil
}

func (p *Proxy) SetStop() {

}
