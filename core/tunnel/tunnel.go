package tunnel

import (
	"net"
	"sync"
	"time"
)

type Option struct {
	Src        net.Conn
	SrcReadLen func(int)
	Dst        net.Conn
	DstReadLen func(int)
	BufLen     int
	Timeout    time.Duration
}

func Tunnel(option Option) {
	if option.Src == nil || option.Dst == nil || option.BufLen <= 1 {
		return
	}
	wg := sync.WaitGroup{}

	fn := func(src, dst net.Conn, readLen func(int)) {
		defer wg.Done()
		defer src.Close()
		for {
			buf := make([]byte, option.BufLen)
			if option.Timeout > 0 {
				src.SetReadDeadline(time.Now().Add(option.Timeout))
			}
			n, err := src.Read(buf)
			if err != nil {
				return
			}
			if readLen != nil {
				go readLen(n)
			}
			if option.Timeout > 0 {
				dst.SetWriteDeadline(time.Now().Add(option.Timeout))
			}
			if wn, err := dst.Write(buf[:n]); err != nil || wn != n {
				return
			}
		}
	}
	wg.Add(2)
	go fn(option.Src, option.Dst, option.SrcReadLen)
	go fn(option.Dst, option.Src, option.DstReadLen)
	wg.Wait()
}
