//go:build !windows

package main

import "errors"

// StartSystem 在非 Windows 平台上不支持，返回错误。
func (p *Proxy) systemStart() error {
	return errors.New("StartSystem is only supported on Windows")
}

// StopSystem 在非 Windows 平台上不支持，返回错误。
func (p *Proxy) systemStop() error {
	return errors.New("StopSystem is only supported on Windows")
}
