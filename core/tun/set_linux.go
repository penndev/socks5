// go:build linux
package tun

func SetAddress(tunName, addr, mask string) error {
	return nil
}

func SetFilterIP(action, addr, mask string) error {
	return nil
}
