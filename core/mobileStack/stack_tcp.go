//go:build linux

package mobileStack

import (
	"fmt"

	"github.com/penndev/socks5/core/socks5"
	"github.com/penndev/socks5/core/stack"
	"github.com/penndev/socks5/core/tunnel"
)

func HandleSocks5TCP() func(ftr *stack.ForwarderTCPRequest) {
	return func(ftr *stack.ForwarderTCPRequest) {
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
				fmt.Println(err)
			}
			return
		}
		option := tunnel.Option{
			Src:    ftr.Conn,
			Dst:    remoteConn,
			BufLen: 32 * 1024,
		}
		if s.Handle != nil {
			if s.Handle.WriteLen != nil {
				option.SrcReadLen = s.Handle.WriteLen
			}
			if s.Handle.ReadLen != nil {
				option.DstReadLen = s.Handle.ReadLen
			}
		}
		tunnel.Tunnel(option)
	}
}

func HandleSocks5UDP() func(fur *stack.ForwarderUDPRequest){
	return func(fur *stack.ForwarderUDPRequest) {
		defer fur.Conn.Close()
		s5, err := socks5.NewClient(srvAddr, s.User, s.Pass)
		if err != nil {
			if s.Handle != nil && s.Handle.Error != nil {
				go s.Handle.Error(err.Error())
			}
			return
		}

		remoteConn, err := s5.Dial("udp", fur.RemoteAddr)
		if err != nil {
			fmt.Println(err)
		}
		option := tunnel.Option{
			Src:     fur.Conn,
			Dst:     remoteConn,
			BufLen:  32 * 1024,
			Timeout: 30 * time.Second,
		}
		if s.Handle != nil {
			if s.Handle.WriteLen != nil {
				option.SrcReadLen = s.Handle.WriteLen
			}
			if s.Handle.ReadLen != nil {
				option.DstReadLen = s.Handle.ReadLen
			}
		}
		tunnel.Tunnel(option)
	}
}