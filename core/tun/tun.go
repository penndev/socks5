package tun

import (
	"context"
	"log"
	"sync"

	"golang.zx2c4.com/wireguard/tun"
	"gvisor.dev/gvisor/pkg/buffer"
	"gvisor.dev/gvisor/pkg/tcpip/header"
	"gvisor.dev/gvisor/pkg/tcpip/link/channel"
	"gvisor.dev/gvisor/pkg/tcpip/stack"
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
	log.Println("i am close")
	(*t.dev).Close()
	t.Endpoint.Close()
}

func (t *Tun) Wait() {
	t.WaitGroup.Wait()
}

func (t *Tun) read(packet []byte) (int, error) {
	bufs := make([][]byte, 1)
	bufs[0] = packet
	sizes := make([]int, 1)
	_, err := (*t.dev).Read(bufs, sizes, 0)
	return sizes[0], err
}

func (t *Tun) write(packet []byte) (int, error) {
	bufs := make([][]byte, 1)
	bufs[0] = packet
	return (*t.dev).Write(bufs, 0)
}

// 从设备读取数据包，并注入到协议栈中
func (t *Tun) inbound(cancel context.CancelFunc) {
	defer t.Done()
	defer cancel()

	for {
		data := make([]byte, int(t.mtu))
		n, err := t.read(data)
		if err != nil {
			// debug
			break
		}
		if n == 0 || n > int(t.mtu) || !t.IsAttached() {
			continue
		}
		pkt := stack.NewPacketBuffer(stack.PacketBufferOptions{
			Payload: buffer.MakeWithData(data),
		})
		switch header.IPVersion(data) {
		case header.IPv4Version:
			t.InjectInbound(header.IPv4ProtocolNumber, pkt)
		case header.IPv6Version:
			t.InjectInbound(header.IPv6ProtocolNumber, pkt)
		}
		pkt.DecRef()
	}
}

// gvisor读取协议栈中数据包，并写入设备
func (t *Tun) outbound(ctx context.Context) {
	defer t.Done()
	for {
		pkt := t.ReadContext(ctx)
		if pkt == nil {
			break
		}
		buf := pkt.ToBuffer()
		t.write(buf.Flatten())
		buf.Release()
		pkt.DecRef()
	}
}

// 协议栈启动
func (t *Tun) Attach(dispatcher stack.NetworkDispatcher) {
	t.Endpoint.Attach(dispatcher)
	t.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		t.Add(2)
		go t.outbound(ctx)
		go t.inbound(cancel)
	})
}

// return stack.LinkEndpoint interface
func CreateTUN(ifname string, mtu int) (*Tun, error) {
	dev, err := tun.CreateTUN(ifname, mtu)
	if err != nil {
		return nil, err
	}
	mtu, err = dev.MTU()
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
