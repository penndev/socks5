package tun

import (
	"context"
	"sync"

	"golang.zx2c4.com/wireguard/tun"
	"gvisor.dev/gvisor/pkg/buffer"
	"gvisor.dev/gvisor/pkg/tcpip/header"
	"gvisor.dev/gvisor/pkg/tcpip/link/channel"
	"gvisor.dev/gvisor/pkg/tcpip/stack"
)

type Tun struct {
	*channel.Endpoint
	mtu    uint32
	offset int
	dev    *tun.Device
	once   sync.Once
	wg     sync.WaitGroup
}

func (t *Tun) Read(packet []byte) (int, error) {
	bufs := make([][]byte, 1)
	bufs[0] = packet
	sizes := make([]int, 1)
	_, err := (*t.dev).Read(bufs, sizes, 0)
	return sizes[0], err
}

func (t *Tun) Write(packet []byte) (int, error) {
	bufs := make([][]byte, 1)
	bufs[0] = packet
	return (*t.dev).Write(bufs, 0)
}

func (t *Tun) Name() string {
	name, _ := (*t.dev).Name()
	return name
}

func (t *Tun) Close() error {
	defer t.Endpoint.Close()
	return (*t.dev).Close()
}

func (t *Tun) Wait() {
	t.wg.Wait()
}

func (t *Tun) Attach(dispatcher stack.NetworkDispatcher) {
	t.Endpoint.Attach(dispatcher)
	t.once.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		t.wg.Add(2)
		go t.readStack(ctx)
		go t.readDev(cancel)
	})
}

func (t *Tun) readStack(ctx context.Context) {
	defer t.wg.Done()
	for {
		pkt := t.ReadContext(ctx)
		if pkt.IsNil() {
			break
		}

		buf := pkt.ToBuffer()
		if t.offset != 0 {
			v := buffer.NewViewWithData(make([]byte, t.offset))
			_ = buf.Prepend(v)
		}
		t.Write(buf.Flatten())
		buf.Release()
		pkt.DecRef()
	}
}

func (t *Tun) readDev(cancel context.CancelFunc) {
	defer t.wg.Done()
	defer cancel()

	for {
		data := make([]byte, t.offset+int(t.mtu))
		n, err := t.Read(data)
		if err != nil {
			// debug
			break
		}
		if n == 0 || n > int(t.mtu) || !t.IsAttached() {
			continue
		}
		pkt := stack.NewPacketBuffer(stack.PacketBufferOptions{
			Payload: buffer.MakeWithData(data[t.offset : t.offset+n]),
		})
		switch header.IPVersion(data[t.offset:]) {
		case header.IPv4Version:
			t.InjectInbound(header.IPv4ProtocolNumber, pkt)
		case header.IPv6Version:
			t.InjectInbound(header.IPv6ProtocolNumber, pkt)
		}
		pkt.DecRef()
	}
}

// return stack.LinkEndpoint interface
func CreateTUN(ifname string, mtu int) (*Tun, error) {
	offset := 0

	dev, err := tun.CreateTUN(ifname, mtu)
	if err != nil {
		return nil, err
	}
	mtu, err = dev.MTU()
	if err != nil {
		return nil, err
	}
	mtu32 := uint32(mtu)
	ep := &Tun{
		mtu:      mtu32,
		dev:      &dev,
		offset:   offset,
		Endpoint: channel.New(1024, mtu32, ""),
	}
	return ep, nil
}
