package tunnel_test

import (
	"fmt"
	"net"
	"time"

	"github.com/penndev/socks5/core/tunnel"
)

func ExampleTunnel() {
	remoteConn, err := net.Dial("tcp", "www.baidu.com:80")
	if err != nil {
		panic(err)
	}

	lr, lw := net.Pipe()

	go tunnel.Tunnel(tunnel.Option{
		Src:     remoteConn,
		Dst:     lw,
		BufLen:  1024,
		Timeout: 10 * time.Second,
	})

	if _, err = lr.Write([]byte("get / \r\n")); err != nil {
		panic(err)
	}
	buf := make([]byte, 1024)
	n, err := lr.Read(buf)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf[:n]))
	// Output: HTTP/1.1 400 Bad Request
}
