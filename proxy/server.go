package proxy

import (
	"log"
	"net"

	"github.com/penndev/gopkg/socks5"
	"github.com/penndev/prism/transport"
)

type Server struct {
	Addr          string
	HandleConnect transport.HandleConnect
	Username      string
	Password      string
	socks5Server  *socks5.Server
}

func (s *Server) handleConn(conn *Conn) {
	// 判断协议类型如果是05则代理socks5，否则代理http
	buf, err := conn.Peek(1)
	if err != nil {
		log.Println("read failed: ", err)
		return
	}
	if buf[0] == 0x05 {
		s.ProxySocks5(conn)
		return
	}
	s.ProxyHTTP(conn)
}

func (s *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept failed: ", err)
			continue
		}
		go s.handleConn(NewConn(conn))
	}
}

func (s *Server) initSocks5() {
	s.socks5Server = &socks5.Server{
		Addr:     s.Addr,
		Username: s.Username,
		Password: s.Password,
		Method:   socks5.METHOD_NO_AUTH,
		HandleConnect: func(c net.Conn, req socks5.Requests, rep socks5.HandleReply) error {
			var err error
			host := req.Addr()
			network := ""
			switch req.CMD {
			case socks5.CMD_CONNECT:
				network = "tcp"
			case socks5.CMD_UDP_ASSOCIATE:
				network = "udp"
			default:
				rep(socks5.REP_COMMAND_NOT_SUPPORTED)
			}
			rep(socks5.REP_SUCCEEDED)
			err = s.HandleConnect(c, network, host)
			return err
		},
	}
	go func() {
		s.socks5Server.UDPListen()
	}()
}

func New(addr, username, password string) *Server {
	s := &Server{
		Addr:          addr,
		HandleConnect: transport.Local(),
		Username:      username,
		Password:      password,
	}
	s.initSocks5()
	return s
}
