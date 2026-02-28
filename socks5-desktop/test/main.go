package main

import (
	"crypto/tls"
	"log"
	"net"

	"github.com/penndev/gopkg/socks5"
)

func HandleConnect(conn net.Conn, req socks5.Requests, replies func(status socks5.REP) error) error {

	addr := req.Addr()
	log.Println("socks5 收到连接 ->", addr)
	// log.Println("logProxyList", "incoming: "+addr)

	var (
		server net.Conn
		err    error
	)

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
	}
	server, err = tls.Dial("tcp", "22cebfa0-tbfog0-tczdfv-1uv4s.sj.gotochinatown.net:443", tlsConf)

	if err != nil {
		log.Println("logProxyList", "connect upstream failed: "+err.Error())
	}

	if err != nil {
		replies(socks5.REP_CONNECTION_REFUSED)
		return err
	}

	socks := &socks5.Client{
		Username: "4ce0c0da-cfed-11f0-86e8-f23c913c8d2b",
		Password: "4ce0c0da-cfed-11f0-86e8-f23c913c8d2b",
		Conn:     server,
	}
	// socks5 握手
	err = socks.Negotiation()
	if err != nil {
		log.Println("logProxyList", "socks5 negotiation failed: "+err.Error())
	}
	if err != nil {
		replies(socks5.REP_CONNECTION_REFUSED)
		return err
	}

	// log.Println("socks5 握手成功", addr)
	// 发起真实的请求
	remote, err := socks.Dial("tcp", addr)
	if err != nil {
		log.Println("logProxyList", "dial target failed: "+err.Error())
		replies(socks5.REP_CONNECTION_REFUSED)
		return err
	}

	replies(socks5.REP_SUCCEEDED)
	defer remote.Close()
	socks5.Pipe(conn, remote)
	// log.Println("logProxyList", "proxy finished: "+addr)
	return nil
}

func main() {
	socks5Server := &socks5.Server{
		Addr:          "127.0.0.1:1080",
		Username:      "",
		Password:      "",
		HandleConnect: HandleConnect,
	}

	log.Println("socks5 启动成功")
	if err := socks5Server.TCPListen(); err != nil {
		log.Println("logServerStatus", "local server start failed: "+err.Error()+"\n")
	}
}
