package tun

import (
	"sync"

	"golang.zx2c4.com/wireguard/tun"
	"gvisor.dev/gvisor/pkg/tcpip/link/channel"
)

type Tun struct {
	*channel.Endpoint
	sync.Once
	sync.WaitGroup

	mtu   uint32
	dev   *tun.Device
	devRM sync.Mutex
	devWM sync.Mutex
}

func (t *Tun) Name() string {
	name, _ := (*t.dev).Name()
	return name
}

func (t *Tun) Close() {
	(*t.dev).Close()
	t.Endpoint.Close()
}

func (t *Tun) Wait() {
	t.WaitGroup.Wait()
}

// return stack.LinkEndpoint interface
func New(options Options) (*Tun, error) {
	dev, err := tun.CreateTUN(options.Name, options.MTU)
	if err != nil {
		return nil, err
	}
	mtu, err := dev.MTU()
	if err != nil {
		return nil, err
	}

	ep := &Tun{
		mtu:      uint32(mtu),
		dev:      &dev,
		Endpoint: channel.New(1024, uint32(mtu), ""),
	}

	return ep, nil
}
