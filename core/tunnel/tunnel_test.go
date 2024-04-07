package tunnel_test

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"testing"
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

func TestTunnel(t *testing.T) {

	rr, rw := net.Pipe()
	lr, lw := net.Pipe()

	go func() {
		for {
			buf := make([]byte, 65535)
			n, err := rr.Read(buf)
			if err != nil {
				return
			}

			if _, e := rr.Write(buf[:n]); e != nil {
				return
			}
		}
	}()

	go func() {
		finish := false
		timeout := time.After(1 * time.Second)
		for {
			select {
			case <-timeout:
				finish = true
			default:
				bufLen := rand.Intn(65535)
				buf := make([]byte, bufLen)
				n, err := lr.Write(buf)
				if err != nil {
					panic(err)
				}
				if n != bufLen {
					panic("n != bufLen")
				}
			}
			if finish {
				break
			}
		}
	}()

	totalReadLen := 0
	go func() {
		for {
			bufRead := make([]byte, 65535)
			n, err := lr.Read(bufRead)
			if err != nil {
				return
			}
			totalReadLen += n
		}
	}()

	option := tunnel.Option{
		Src:     rw,
		Dst:     lw,
		BufLen:  1024,
		Timeout: 1 * time.Second,
	}

	var mu sync.Mutex
	srl := 0
	option.SrcReadLen = func(i int) {
		mu.Lock()
		defer mu.Unlock()
		srl += i
	}

	var wmu sync.Mutex
	drl := 0
	option.DstReadLen = func(i int) {
		wmu.Lock()
		defer wmu.Unlock()
		drl += i
	}

	tunnel.Tunnel(option)

	log.Println("MBps/", totalReadLen*8/1000000, "|", srl, "|", drl)
}
