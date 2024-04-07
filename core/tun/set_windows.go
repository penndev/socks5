// go:build windows
package tun

import (
	"fmt"
	"os/exec"
	"strings"
)

// SetAddress is set tun static ip
func SetAddress(tunName, addr, mask string) error {
	arg := fmt.Sprintf(`interface ipv4 set address name="%s" source=static addr=%s mask=%s`, tunName, addr, mask)
	return exec.Command("netsh", strings.Fields(arg)...).Run()
}

// action
// - route add ...
// - route delete ...
func SetFilterIP(action, addr, mask string) error {
	arg := fmt.Sprintf(` %s mask %s 0.0.0.0 metric 1`, addr, mask)
	arg = action + arg
	return exec.Command("route", strings.Fields(arg)...).Run()
}
