package socks5

import (
	"io"
	"sync"
)

func TunnelTCP(dst, src io.ReadWriteCloser) {
	var wg sync.WaitGroup
	wg.Add(2)
	fn := func(dst, src io.ReadWriteCloser) {
		io.Copy(dst, src)
		wg.Done()
	}
	go fn(dst, src)
	go fn(src, dst)
	wg.Wait()
}
