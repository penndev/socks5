//go:build linux

package mobileStack

import (
	"fmt"

	"github.com/penndev/socks5/core/fdtun"
	"github.com/penndev/socks5/core/socks5"
	"github.com/penndev/socks5/core/stack"
)

func Version() string {
	return "0.0.1"
}

type StackHandle interface {
	Error(s string)
	WriteByteLen(l int)
	ReadByteLen(l int)
}

type Stack struct {
	TunFd   int
	MTU     int
	User    string
	Pass    string
	SrvHost string
	SrvPort int
	Handle  StackHandle
}

func (s *Stack) Run() (bool, error) {
	if s.TunFd < 1 {
		return false, fmt.Errorf("tunFd < 1:[%d]")
	}
	if s.MTU < 64 {
		return false, fmt.Errorf("mtu < 64:[%d]")
	}
	if len(s.SrvHost) < 4 {
		return false, fmt.Errorf("srvHost < 4:[%d]")
	}
	if s.SrvPort < 1 {
		return false, fmt.Errorf("srvPort < 1:[%d]")
	}
	dev, err := fdtun.CreateTUN(s.TunFd, uint32(s.MTU))
	if err != nil {
		return false, err
	}
	srvAddr := fmt.Sprintf("%s:%d", s.SrvHost, s.SrvPort)
	stack.New(stack.Option{
		EndPoint: dev,
		HandleTCP: func(ftr *stack.ForwarderTCPRequest) {
			defer ftr.Conn.Close()
			s5, err := socks5.NewClient(srvAddr, s.User, s.Pass)
			if err != nil {
				if s.Handle != nil && s.Handle.Error != nil {
					go s.Handle.Error(err.Error())
				}
				return
			}
			defer s5.Close()

			remoteConn, err := s5.Dial("tcp", ftr.RemoteAddr)
			if err != nil {
				if s.Handle != nil && s.Handle.Error != nil {
					go s.Handle.Error(err.Error())
				}
				return
			}
			socks5.TunnelTCP(ftr.Conn, remoteConn)
		},
	})
	return true, nil
}

func NewStack() *Stack { return &Stack{} }
