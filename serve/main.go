package main

import (
	"flag"
	"log"
	"net"

	"github.com/penndev/prism/proxy"
	"github.com/penndev/prism/transport"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "listen address")
	user := flag.String("user", "", "proxy username")
	pass := flag.String("pass", "", "proxy password")
	flag.Parse()

	s := proxy.New(*addr, *user, *pass)
	handle := transport.Local()
	s.HandleConnect = func(conn net.Conn, network, address string) error {
		log.Println("req ->", network, address)
		return handle(conn, network, address)
	}
	log.Printf("start -> %s %s %s", *addr, *user, *pass)
	err := s.ListenAndServe()
	if err != nil {
		log.Println("listen failed: ", err)
	}
}
