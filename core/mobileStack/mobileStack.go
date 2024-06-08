//go:build linux

package mobileStack

import (
	"fmt"

	"github.com/penndev/socks5/core/fdtun"
	"github.com/penndev/socks5/core/stack"
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
	TunFd int //tun设置代理
	MTU   int //mtu掩码

	SrvHost string
	SrvPort int
	User    string
	Pass    string

	TcpEnable bool        //是否启用tcp代理
	UdpEnable bool        //是否启用udp代理
	Handle    StackHandle //各种事件回调
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

	stackOption := stack.Option{
		EndPoint: dev,
	}

	stackOption.HandleTCP = HandleSocks5TCP(s)
	stackOption.HandlerUDP = HandleSocks5UDP(s)
	stack.New(stackOption)
	return true, nil
}

func NewStack() *Stack { return &Stack{} }
