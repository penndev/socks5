package route

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/netip"
	"os/exec"
	"strings"
	"time"
)

func prefixMask(p netip.Prefix) net.IP {
	if !p.IsValid() {
		return nil
	}
	if p.Addr().Is4() {
		return net.IP(net.CIDRMask(p.Bits(), 32))
	}
	return net.IP(net.CIDRMask(p.Bits(), 128))
}

func waitRouteReady(tunPrefix netip.Prefix) error {

	tunIP := tunPrefix.Addr().String()
	ticker := time.NewTicker(200 * time.Millisecond)
	timeout := time.After(30 * time.Second)

	for {
		select {
		case <-ticker.C:
			out, err := exec.Command("ipconfig").CombinedOutput()
			if err != nil {
				continue
			}
			if strings.Contains(string(out), tunIP) {
				return nil
			} else {
				continue
			}
		case <-timeout:
			return errors.New("cant set dev static ip")
		}
	}
}

// 给设备设置静态IP
func SetDevAddr(tunName string, prefix netip.Prefix) error {
	if !prefix.Addr().Is4() {
		panic("not ipv4 prefix")
	}

	args := []string{
		"interface", "ipv4", "set", "address",
		fmt.Sprintf(`name=%s`, tunName),
		"source=static",
		fmt.Sprintf("addr=%s", prefix.Addr().String()),
		fmt.Sprintf("mask=%s", prefixMask(prefix).String()),
	}
	log.Println("netsh", strings.Join(args, " "))

	out, err := exec.Command("netsh", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("netsh failed: %v, %s", err, string(out))
	}
	return waitRouteReady(prefix)
}

// 设置路由表
// 设置路由表
func SetRouteAddr(addr netip.Prefix, gateway net.IP) error {
	if !addr.IsValid() {
		return fmt.Errorf("invalid route prefix")
	}
	if gateway == nil {
		return fmt.Errorf("gateway is nil")
	}

	args := []string{
		"add",
		addr.Addr().String(),
		"mask",
		prefixMask(addr).String(),
		gateway.String(),
	}

	log.Println("route", strings.Join(args, " "))

	out, err := exec.Command("route", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("route failed: %v, output: %s", err, string(out))
	}

	if len(out) > 0 {
		log.Println("route output:", string(out))
	}
	return nil
}
