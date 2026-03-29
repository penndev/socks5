package main

import (
	"fmt"
	"testing"
)

func TestProxyPing(t *testing.T) {
	proxyPing := &ProxyPing{}
	// result := proxyPing.TestServer("socks5overtls://4ce0c0da-cfed-11f0-86e8-f23c913c8d2b:4ce0c0da-cfed-11f0-86e8-f23c913c8d2b@ae73b657-t9hb40-tczdfv-1uv4s.sjc.oshuawei.com:443")
	result := proxyPing.TestServer("socks5://127.0.0.1:1080")
	fmt.Println(result)
	t.Fail()
}
