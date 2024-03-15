package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/penndev/socks5/core/socks5"
	"github.com/penndev/socks5/core/stack"
	"github.com/penndev/socks5/core/tun"
)

func main() {
	tun, err := tun.CreateTUN("wintun", 0)
	if err != nil {
		panic(err)
	}

	stack.Start(stack.Option{
		EndPoint: tun,
		HandleTCP: func(ftr *stack.ForwarderTCPRequest) {
			addrInfo := ftr.ID()
			log.Printf("local %s:%d <- ... -> %s:%d remote",
				addrInfo.LocalAddress, addrInfo.LocalPort, addrInfo.RemoteAddress, addrInfo.RemotePort)

			localConn, err := ftr.TCPConn()
			if err != nil {
				panic(err)
			}
			defer localConn.Close()

			s5, err := socks5.NewClient("127.0.0.1:1080", "", "")
			if err != nil {
				panic(err)
			}
			reqAddr := fmt.Sprintf("%s:%d", addrInfo.LocalAddress, addrInfo.LocalPort)
			remoteConn, err := s5.Dial("tcp", reqAddr)
			if err != nil {
				panic(err)
			}
			defer remoteConn.Close()
			socks5.TunnelTCP(localConn, remoteConn)
		},
	})

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
