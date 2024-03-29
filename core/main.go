// go:build windows
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/penndev/socks5/core/socks5"
	"github.com/penndev/socks5/core/stack"
	"github.com/penndev/socks5/core/tun"
)

func main() {
	dev, err := tun.CreateTUN("wintun", 0)
	if err != nil {
		panic(err)
	}
	if err := tun.Cfg(dev.Name(), "10.10.100.251", "255.255.255.255"); err != nil {
		panic("cant set ip")
	}
	stack.New(stack.Option{
		EndPoint: dev,
		HandleTCP: func(ftr *stack.ForwarderTCPRequest) {
			defer ftr.Conn.Close()
			s5, err := socks5.NewClient("127.0.0.1:1080", "", "")
			if err != nil {
				log.Println("socks5 connection err:", err)
				return
			}
			defer s5.Close()

			remoteConn, err := s5.Dial("tcp", ftr.RemoteAddr)
			if err != nil {
				log.Println("socks5 remote err:", err)
				return
			}
			socks5.TunnelTCP(ftr.Conn, remoteConn)
		},
	})

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
