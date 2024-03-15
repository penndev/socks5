package socks5_test

import (
	"fmt"

	"github.com/penndev/socks5/core/socks5"
)

func ExampleNewClient() {
	s5, err := socks5.NewClient("127.0.0.1:1080", "", "")
	if err != nil {
		panic(err)
	}
	conn, err := s5.Dial("tcp", "www.baidu.com:80")
	if err != nil {
		panic(err)
	}
	_, err = conn.Write([]byte("get / \r\n"))
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 102400)
	n, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf[:n]))
	// Output: HTTP/1.1 400 Bad Request
}
