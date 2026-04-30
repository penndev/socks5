package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"
	"net/url"

	"github.com/penndev/prism/proxy"
	"github.com/penndev/prism/transport"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:1080", "listen address")
	user := flag.String("user", "", "proxy username")
	pass := flag.String("pass", "", "proxy password")
	proxyurl := flag.String("proxy", "", "Set remote socks5 service example socks5://user:pass@192.168.0.1:1080")
	flag.Parse()

	handle := transport.Local()

	if *proxyurl != "" {
		proxyURL, err := url.Parse(*proxyurl)
		if err != nil {
			panic(err)
		}
		username := proxyURL.User.Username()
		password, _ := proxyURL.User.Password()
		switch proxyURL.Scheme {
		case "socks5":
			handle = transport.Socks5(
				proxyURL.Host,
				username,
				password,
			)
		case "socks5tls":
			handle = transport.Socks5OverTLS(
				proxyURL.Host,
				username,
				password,
				&tls.Config{},
			)
		case "http":
			handle = transport.Http(
				proxyURL.Host,
				username,
				password,
			)
		case "https":
			handle = transport.HttpOverTLS(
				proxyURL.Host,
				username,
				password,
				&tls.Config{},
			)
		}
	}
	s := proxy.New(*addr, *user, *pass)
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
