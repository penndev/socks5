//go:build linux

package mobileStack

import (
	"fmt"
	"log"

	"github.com/penndev/socks5/core/fdtun"
	"github.com/penndev/socks5/core/socks5"
	"github.com/penndev/socks5/core/stack"
)

type Option struct {
	TunFd       int
	MTU         uint32
	User        string
	Pass        string
	SrvHost     string
	SrvPort     uint16
	HandleError func(string)
}

func New(option Option) error {
	dev, err := fdtun.CreateTUN(option.TunFd, option.MTU)
	if err != nil {
		return err
	}
	srvAddr := fmt.Sprintf("%s:%d", option.SrvHost, option.SrvPort)
	stack.New(stack.Option{
		EndPoint: dev,
		HandleTCP: func(ftr *stack.ForwarderTCPRequest) {
			defer ftr.Conn.Close()
			s5, err := socks5.NewClient(srvAddr, option.User, option.Pass)
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
	return nil
}
