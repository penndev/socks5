//go:build linux

package mobileStack

import (
	"fmt"
	"time"

	"github.com/penndev/socks5/core/fdtun"
	"github.com/penndev/socks5/core/socks5"
	"github.com/penndev/socks5/core/stack"
	"github.com/penndev/socks5/core/tunnel"
)

func Version() string {
	return "0.0.1"
}

type StackHandle interface {
	Error(string)
	WriteLen(int)
	ReadLen(int)
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

	stackOption := stack.Option{
		EndPoint: dev,
	}

	stackOption.HandleTCP = HandleSocks5TCP()
	stackOption.HandlerUDP = HandleSocks5UDP()
	stack.New(stackOption)
	return true, nil
}

func NewStack() *Stack { return &Stack{} }
