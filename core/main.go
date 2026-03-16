package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/penndev/socks5/core/netlink"
	"github.com/penndev/socks5/core/stack"
	"github.com/penndev/socks5/core/tun"
)

type Config struct {
	Proxy string
}

var (
	config = Config{}
)

func init() {
	flag.StringVar(&config.Proxy, "proxy", "", "Set remote socks5 service example socks5://user:pass@192.168.0.1:1080")
	flag.Parse()
}

func main() {

	u, err := url.Parse(config.Proxy)
	if err != nil {
		panic(err)
	}
	if u.Scheme != "socks5" {
		panic("scheme error")
	}
	dev, err := tun.CreateTUN(TUN_NAME, 0)
	if err != nil {
		panic(err)
	}
	err = netlink.SetAddress(TUN_NAME, "172.19.0.1", "255.255.255.255")
	if err != nil {
		panic(err)
	}
	defer dev.Close()
	stack.New(stack.Option{
		EndPoint: dev,
		HandleTCP: func(ftr *stack.ForwarderTCPRequest) {
			defer ftr.Conn.Close()
			log.Printf("tcp %s <-> %s", ftr.LocalAddr, ftr.RemoteAddr)
		},
	})

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
