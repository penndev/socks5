package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/netip"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/penndev/socks5/core/internal"
	"github.com/penndev/socks5/core/internal/route"
	"github.com/penndev/socks5/core/stack"
	"github.com/penndev/socks5/core/tun"
)

var proxy string

func init() {
	flag.StringVar(&proxy, "proxy", "", "Set remote socks5 service example socks5://user:pass@192.168.0.1:1080")
	flag.Parse()
}

func main() {
	var handleConnect internal.HandleConnect
	handleConnect = internal.Local()
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		panic(err)
	}
	username := proxyURL.User.Username()
	password, _ := proxyURL.User.Password()
	switch proxyURL.Scheme {
	case "socks5":
		handleConnect = internal.Socks5(
			proxyURL.Host,
			username,
			password,
		)
	case "socks5tls":
		handleConnect = internal.Socks5OverTLS(
			proxyURL.Host,
			username,
			password,
			&tls.Config{},
		)
	}

	dev, err := tun.New(tun.Options{
		Name:   TUN_NAME,
		MTU:    TUN_MTU,
		Offset: TUN_OFFSET,
	})
	if err != nil {
		panic(err)
	}
	defer dev.Close()
	stack.New(stack.Option{
		EndPoint: dev,
		HandleTCP: func(f *stack.ForwarderTCPRequest) {
			log.Printf(
				"%s -> %s",
				f.RemoteAddr.Network(),
				f.RemoteAddr.String(),
			)
			handleConnect(f.Conn, &f.RemoteAddr)
		},
		HandlerUDP: func(f *stack.ForwarderUDPRequest) {
			log.Printf(
				"%s -> %s",
				f.RemoteAddr.Network(),
				f.RemoteAddr.String(),
			)
			handleConnect(f.Conn, &f.RemoteAddr)
		},
	})
	route.Start(route.Options{
		DevName:      dev.Name(),
		DevIP:        netip.MustParsePrefix("172.19.0.1/32"),
		RouteAddress: []netip.Prefix{netip.MustParsePrefix("0.0.0.0/0")},
	})

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
