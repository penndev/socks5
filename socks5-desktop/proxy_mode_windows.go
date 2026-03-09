//go:build windows

package main

import (
	"errors"

	"golang.org/x/sys/windows/registry"
)

// StartSystem 设置 Windows 系统代理为本地 socks5 服务地址（p.localServer.Addr）。
func (p *Proxy) systemStart() error {
	if p.localServer == nil || p.localServer.Addr == "" {
		return errors.New("local socks5 server not started")
	}

	addr := p.localServer.Addr
	// Windows 下 socks5 代理格式：socks=ip:port
	proxyValue := "socks=" + addr

	key, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Internet Settings`,
		registry.SET_VALUE,
	)
	if err != nil {
		return err
	}
	defer key.Close()

	if err := key.SetStringValue("ProxyServer", proxyValue); err != nil {
		return err
	}
	if err := key.SetDWordValue("ProxyEnable", 1); err != nil {
		return err
	}

	return nil
}

// StopSystem 关闭 Windows 系统代理（不再通过本地 socks5）。
func (p *Proxy) systemStop() error {
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Internet Settings`,
		registry.SET_VALUE,
	)
	if err != nil {
		return err
	}
	defer key.Close()

	if err := key.SetDWordValue("ProxyEnable", 0); err != nil {
		return err
	}

	return nil
}
