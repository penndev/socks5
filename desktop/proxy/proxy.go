package proxy

import (
	"net/netip"
	"net/url"

	"desktop/internal"

	"github.com/penndev/prism/proxy"
	"github.com/penndev/prism/route"
	"github.com/penndev/prism/stack"
	"github.com/penndev/prism/tun"
)

type Proxy struct {
	proxy.Server

	// tun用
	dev *tun.Tun
}

func (p *Proxy) SetStart(host, user, pass string) error {
	internal.App.Event.Emit(
		internal.AppConfig.LogTypeName_STATUS,
		"localServer://"+user+":"+pass+"@"+host,
	)

	// 配置未变化，保持当前服务
	if p.Addr == host && p.Username == user && p.Password == pass {
		return nil
	}

	p.Server.Close()
	go func() {
		p.Addr = host
		p.Username = user
		p.Password = pass
		if err := p.ListenAndServe(); err != nil {
			internal.App.Event.Emit(
				internal.AppConfig.LogTypeName_STATUS,
				"p.ListenAndServe error: "+err.Error(),
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
	p.HandleConnect, err = HandleConnect(ru)
	if err != nil {
		internal.App.Event.Emit(internal.AppConfig.LogTypeName_LOG, "SetRemote error: "+err.Error())
		return err
	}
	return nil
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
	stack.New(stack.Option{
		EndPoint: p.dev,
		HandleTCP: func(f *stack.ForwarderTCPRequest) {
			internal.App.Event.Emit(internal.AppConfig.LogTypeName_LOG, "tun -> "+f.RemoteAddr.Network()+" "+f.RemoteAddr.String())
			p.HandleConnect(f.Conn, f.RemoteAddr.Network(), f.RemoteAddr.String())
		},
		HandlerUDP: func(f *stack.ForwarderUDPRequest) {
			internal.App.Event.Emit(internal.AppConfig.LogTypeName_LOG, "tun -> "+f.RemoteAddr.Network()+" "+f.RemoteAddr.String())
			p.HandleConnect(f.Conn, f.RemoteAddr.Network(), f.RemoteAddr.String())
		},
	})
	route.Start(route.Options{
		DevName:      p.dev.Name(),
		DevIP:        netip.MustParsePrefix("172.19.0.1/32"),
		RouteAddress: []netip.Prefix{netip.MustParsePrefix("0.0.0.0/0")},
	})
	return nil
}

func (p *Proxy) SetMode(mode string) {
	internal.App.Event.Emit(
		internal.AppConfig.LogTypeName_STATUS,
		"SetMode-> "+mode,
	)
	var err error
	switch mode {
	case "tun":
		err = p.setModeTun()
	default:
		if p.dev != nil {
			p.dev.Close()
			internal.App.Event.Emit(internal.AppConfig.LogTypeName_STATUS, "tun dev close success")
		}
	}
	if err != nil {
		internal.App.Event.Emit(internal.AppConfig.LogTypeName_STATUS, err.Error())
	}
}

func (p *Proxy) SetStop() {
	if p.dev != nil {
		p.dev.Close()
		internal.App.Event.Emit(
			internal.AppConfig.LogTypeName_STATUS,
			"tun dev close success",
		)
	}
	p.Server.Close()
	internal.App.Event.Emit(
		internal.AppConfig.LogTypeName_STATUS,
		"local server close success",
	)
}
