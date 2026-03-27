package main

import (
	"crypto/tls"
	"errors"
	"net"
	"net/url"
	"time"
)

type Ping struct{}

type PingRequest struct {
	Host     string `json:"host"`
	Protocol string `json:"protocol"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type PingResult struct {
	Latency int    `json:"latency"` // 延迟（毫秒）
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// TestServer 测试服务器连接延迟
func (p *Ping) TestServer(req PingRequest) PingResult {
	start := time.Now()

	// 构建远程 URL
	var remoteURL string
	if req.Username != "" && req.Password != "" {
		remoteURL = req.Protocol + "://" + req.Username + ":" + req.Password + "@" + req.Host
	} else {
		remoteURL = req.Protocol + "://" + req.Host
	}

	ru, err := url.Parse(remoteURL)
	if err != nil {
		return PingResult{
			Success: false,
			Error:   "解析地址失败: " + err.Error(),
		}
	}

	// 测试连接
	err = testConnection(ru)
	elapsed := time.Since(start)

	if err != nil {
		return PingResult{
			Success: false,
			Error:   err.Error(),
		}
	}

	return PingResult{
		Success: true,
		Latency: int(elapsed.Milliseconds()),
	}
}

// testConnection 测试到 SOCKS5 服务器的连接
func testConnection(remote *url.URL) error {
	timeout := 5 * time.Second

	switch remote.Scheme {
	case "socks5":
		// 直接 TCP 连接测试
		conn, err := net.DialTimeout("tcp", remote.Host, timeout)
		if err != nil {
			return errors.New("连接失败: " + err.Error())
		}
		defer conn.Close()

		// 简单的 SOCKS5 握手测试
		// 发送版本协商请求
		_, err = conn.Write([]byte{0x05, 0x01, 0x00})
		if err != nil {
			return errors.New("发送握手失败: " + err.Error())
		}

		// 读取服务器响应
		conn.SetReadDeadline(time.Now().Add(timeout))
		buf := make([]byte, 2)
		_, err = conn.Read(buf)
		if err != nil {
			return errors.New("读取响应失败: " + err.Error())
		}

		if buf[0] != 0x05 {
			return errors.New("无效的 SOCKS5 响应")
		}

	case "socks5overtls":
		// TLS 连接测试
		conn, err := tls.DialWithDialer(
			&net.Dialer{Timeout: timeout},
			"tcp",
			remote.Host,
			&tls.Config{InsecureSkipVerify: true},
		)
		if err != nil {
			return errors.New("TLS 连接失败: " + err.Error())
		}
		defer conn.Close()

		// 简单的 SOCKS5 握手测试
		_, err = conn.Write([]byte{0x05, 0x01, 0x00})
		if err != nil {
			return errors.New("发送握手失败: " + err.Error())
		}

		conn.SetReadDeadline(time.Now().Add(timeout))
		buf := make([]byte, 2)
		_, err = conn.Read(buf)
		if err != nil {
			return errors.New("读取响应失败: " + err.Error())
		}

		if buf[0] != 0x05 {
			return errors.New("无效的 SOCKS5 响应")
		}

	default:
		return errors.New("不支持的协议: " + remote.Scheme)
	}

	return nil
}
