package dialer

import (
	"io"
	"net"
)

func Pipe(src, dst net.Conn) {
	go func() {
		defer dst.Close()
		io.Copy(dst, src)
	}()
	defer src.Close()
	io.Copy(src, dst)
}
