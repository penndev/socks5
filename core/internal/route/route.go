package route

import (
	"net/netip"
)

type Options struct {
	DevName      string
	DevIP        netip.Prefix
	RouteAddress []netip.Prefix
}

func Start(options Options) error {
	// 设置路由静态IP与掩码位
	err := SetDevAddr(options.DevName, options.DevIP)
	if err != nil {
		return err
	}
	for _, item := range options.RouteAddress {
		SetRouteAddr(item, options.DevIP.Addr().AsSlice())
	}
	return err
}
