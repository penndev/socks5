// go:build windows
package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/penndev/socks5/core/socks5"
	"github.com/penndev/socks5/core/stack"
	"github.com/penndev/socks5/core/tun"
	"github.com/penndev/socks5/core/tunnel"
)

func main() {
	if versionFlag {
		fmt.Println(version)
		return
	}
	u, err := url.Parse(config.Proxy)
	if err != nil {
		panic(err)
	}
	if u.Scheme != "socks5" {
		panic("scheme error")
	}
	host := u.Host
	user := u.User.Username()
	pass, _ := u.User.Password()

	dev, err := tun.CreateTUN(config.TunName, 0)
	if err != nil {
		panic(err)
	}
	if err := tun.Cfg(dev.Name(), config.TunIP, "255.255.255.255"); err != nil {
		panic("cant set static ip")
	}
	stack.New(stack.Option{
		EndPoint: dev,
		HandleTCP: func(ftr *stack.ForwarderTCPRequest) {
			defer ftr.Conn.Close()
			s5, err := socks5.NewClient(host, user, pass)
			if err != nil {
				log.Println("socks5 connection err:", err)
				return
			}
			defer s5.Close()

			remoteConn, err := s5.Dial("tcp", ftr.RemoteAddr)
			if err != nil {
				log.Println("socks5 tcp remote err:", err)
				return
			}
			log.Printf("tcp %s <-> %s", ftr.LocalAddr, ftr.RemoteAddr)
			tunnel.Tunnel(tunnel.Option{
				Src:    ftr.Conn,
				Dst:    remoteConn,
				BufLen: 32 * 1024,
			})
		},
		HandlerUDP: func(fur *stack.ForwarderUDPRequest) {
			defer fur.Conn.Close()
			s5, err := socks5.NewClient(host, user, pass)
			if err != nil {
				log.Println("socks5 connection err:", err)
				return
			}
			defer s5.Close()

			remoteConn, err := s5.Dial("udp", fur.RemoteAddr)
			if err != nil {
				log.Println("socks5 udp remote err:", err)
				return
			}
			log.Printf("udp %s <-> %s", fur.LocalAddr, fur.RemoteAddr)
			tunnel.Tunnel(tunnel.Option{
				Src:     fur.Conn,
				Dst:     remoteConn,
				BufLen:  1024,
				Timeout: 30 * time.Second,
			})
		},
	})

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
