package proxy

import (
	"net"

	"github.com/penndev/gopkg/socks5"
)

func (s *Server) ProxySocks5(conn net.Conn) {
	if s.socks5Server.Username != "" {
		s.socks5Server.Method = socks5.METHOD_USERNAME_PASSWORD
	}
	s.socks5Server.HandleConn(conn)
}
