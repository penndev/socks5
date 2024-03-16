//go:build linux

package fdtun

import (
	"gvisor.dev/gvisor/pkg/tcpip/link/fdbased"
	"gvisor.dev/gvisor/pkg/tcpip/stack"
)

type FD struct {
	stack.LinkEndpoint
}

func CreateTUN(fd int, mtu uint32) (stack.LinkEndpoint, error) {
	ep, err := fdbased.New(&fdbased.Options{
		FDs:            []int{fd},
		MTU:            mtu,
		EthernetHeader: false,
	})
	if err != nil {
		return nil, err
	}
	dev := &FD{
		LinkEndpoint: ep,
	}
	return dev, nil
}
