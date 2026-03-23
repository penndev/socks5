package dialer

import (
	"io"
	"net"
	"sync"
)

func Pipe(src, dst net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)

	// 对 UDP 这类“请求->响应”场景，src->dst 的 io.Copy 可能会因为单次包而提前结束。
	// 旧实现会在这个时刻提前关闭 dst，导致对方响应尚未读到就被丢弃。
	// 因此这里不在发送方向里关闭 dst，改为在两个方向结束后统一关闭。
	go func() {
		defer wg.Done()
		_, _ = io.Copy(dst, src)
	}()

	go func() {
		defer wg.Done()
		_, _ = io.Copy(src, dst)
	}()

	wg.Wait()
	_ = src.Close()
	_ = dst.Close()
}
