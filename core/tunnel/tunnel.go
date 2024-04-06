package tunnel

import (
	"net"
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

func readChan(ch chan []byte, conn net.Conn, bufLen int, readLen func(int)) {
	defer close(ch)
	if readLen == nil {
		for {
			buf := make([]byte, bufLen)
			n, err := conn.Read(buf)
			if err != nil {
				break
			}
			ch <- buf[:n]
		}
	} else {
		for {
			buf := make([]byte, bufLen)
			n, err := conn.Read(buf)
			readLen(n)
			if err != nil {
				break
			}
			ch <- buf[:n]
		}
	}
}

func Tunnel(option Option) {
	if option.Src == nil || option.Dst == nil || option.BufLen <= 1 {
		return
	}
	srcChan := make(chan []byte)
	go readChan(srcChan, option.Src, option.BufLen, option.SrcReadLen)

	dstChan := make(chan []byte)
	go readChan(dstChan, option.Dst, option.BufLen, option.DstReadLen)

	if option.Timeout > 0 {
		timer := time.NewTimer(option.Timeout)
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				return
			case buf, ok := <-dstChan:
				if !ok {
					return
				}
				if n, err := option.Src.Write(buf); err != nil || n != len(buf) {
					return
				}
			case buf, ok := <-srcChan:
				if !ok {
					return
				}
				if n, err := option.Dst.Write(buf); err != nil || n != len(buf) {
					return
				}
			}
			timer.Reset(option.Timeout)
		}
	} else {
		for {
			select {
			case buf, ok := <-dstChan:
				if !ok {
					return
				}
				if n, err := option.Src.Write(buf); err != nil || n != len(buf) {
					return
				}
			case buf, ok := <-srcChan:
				if !ok {
					return
				}
				if n, err := option.Dst.Write(buf); err != nil || n != len(buf) {
					return
				}
			}
		}
	}
}
