// go:build linux
package tun

import (
	"fmt"
	"os/exec"
	"strings"
)

const netshIPv4Args = `interface ipv4 set address name="%s" source=static addr=%s mask=%s`

func Cfg(tunName, addr, mask string) error {
	arg := fmt.Sprintf(netshIPv4Args, tunName, addr, mask)
	return exec.Command("netsh", strings.Fields(arg)...).Run()
}
