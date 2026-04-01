package proxy

import (
	"net"
	"net/url"
	"time"
)

type ProxyPing struct{}

type ProxyPingResult struct {
	Latency int    `json:"latency"` // 延迟（毫秒）
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// TestServer 通过 HandleConnect 经代理访问 Google（HTTP）并测量首字节响应延迟。
func (p *ProxyPing) TestServer(serverURL string) ProxyPingResult {
	remote, err := url.Parse(serverURL)
	if err != nil {
		return ProxyPingResult{Success: false, Error: err.Error()}
	}
	handle, err := HandleConnect(remote)
	if err != nil {
		return ProxyPingResult{Success: false, Error: err.Error()}
	}

	c1, c2 := net.Pipe()
	defer c1.Close()
	defer c2.Close()

	// 建立proxy连接
	overCH := make(chan error, 1)
	go func() {
		overCH <- handle(c1, "tcp", "www.google.com:80")
	}()

	start := time.Now()
	deadline := 5 * time.Second
	_ = c2.SetDeadline(start.Add(deadline))
	if _, err := c2.Write([]byte("GET / HTTP/1.0\r\nHost: www.google.com\r\n\r\n")); err != nil {
		return ProxyPingResult{Success: false, Error: err.Error()}
	}
	go func() {
		buf := make([]byte, 1)
		_, err := c2.Read(buf)
		overCH <- err
	}()

	timer := time.NewTimer(deadline)
	defer timer.Stop()
	select {
	case err := <-overCH:
		if err != nil {
			return ProxyPingResult{Success: false, Error: err.Error()}
		} else {
			return ProxyPingResult{Success: true, Latency: int(time.Since(start).Milliseconds())}
		}
	case <-timer.C:
		return ProxyPingResult{Success: false, Error: "timeout"}
	}
}
