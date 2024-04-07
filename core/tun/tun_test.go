package tun_test

import (
	"fmt"

	"github.com/penndev/socks5/core/tun"
)

func ExampleSetAddress() {
	fmt.Println(tun.SetAddress("socks5", "10.1.1.1", "255.255.255.255"))
	// output:
	// exit status 1
}
