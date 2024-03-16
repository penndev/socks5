package tun_test

import (
	"fmt"

	"github.com/penndev/socks5/core/tun"
)

func ExampleCfg() {
	fmt.Println(tun.Cfg("wintun", "10.1.1.1", "255.255.255.255"))
	// output:
	// exit status 1
}
