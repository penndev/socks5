package tunnel

import (
	"io"
	"time"
)

func Tunnel(dst, src io.ReadWriteCloser, bufferLen int, timeout time.Duration) {
	fn := func(rw io.ReadWriteCloser, ch chan []byte) {
		defer close(ch)
		for {
			buf := make([]byte, bufferLen)
			n, err := rw.Read(buf)
			if err != nil {
				break
			}
			ch <- buf[:n]
		}
	}

	dstChan := make(chan []byte)
	go fn(dst, dstChan)

	srcChan := make(chan []byte)
	go fn(src, srcChan)

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			return
		case buf, ok := <-dstChan:
			if !ok {
				return
			}
			if n, err := src.Write(buf); err != nil || n != len(buf) {
				return
			}
		case buf, ok := <-srcChan:
			if !ok {
				return
			}
			if n, err := dst.Write(buf); err != nil || n != len(buf) {
				return
			}
		}
		timer.Reset(timeout)
	}
}
