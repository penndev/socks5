package proxy

import (
	"fmt"
	"testing"
)

func TestProxyPing(t *testing.T) {
	proxyPing := &ProxyPing{}
	result := proxyPing.TestServer("socks5://127.0.0.1:1080")
	fmt.Println(result)
	t.Fail()
}
