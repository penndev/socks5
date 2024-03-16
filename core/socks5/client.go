package socks5

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
)

type Conn struct {
	rw net.Conn
}

func NewClient(address, user, pass string) (Conn, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return Conn{}, err
	}
	var connects []byte
	if user == "" {
		connects = []byte{Version, 0x1, METHOD_NO_AUTH}
	} else {
		connects = []byte{Version, 0x2, METHOD_NO_AUTH, METHOD_USER}
	}
	if _, err := conn.Write(connects); err != nil {
		conn.Close()
		return Conn{}, err
	}
	buf := make([]byte, 2)
	rn, err := conn.Read(buf)
	if err != nil {
		conn.Close()
		return Conn{}, err
	}
	if rn != 2 || buf[0] != Version {
		conn.Close()
		return Conn{}, errors.New("error socks5 service Version")
	}
	switch buf[1] {
	case METHOD_NO_AUTH:
		return Conn{rw: conn}, err
	default:
		conn.Close()
		return Conn{}, fmt.Errorf("error socks method not match [%d]", buf[1])
	}
}

func (c *Conn) requests(network, address string) (Requests, error) {
	req := Requests{}

	if network == "tcp" {
		req.CMD = 0x01
	} else {
		return req, errors.New("暂不支持 " + network)
	}

	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return req, err
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return req, err
	} else if portInt < 0 || portInt > 65535 {
		return req, errors.New("port error")
	}
	req.DST_PORT = uint16(portInt)

	dstAddr := net.ParseIP(host)
	if dstAddr == nil {
		req.ATYP = 0x03
		req.DST_ADDR = []byte(host)
	} else if ip4 := dstAddr.To4(); ip4 != nil {
		req.ATYP = 0x01
		req.DST_ADDR = []byte(ip4)
	} else if ip6 := dstAddr.To16(); ip6 != nil {
		req.ATYP = 0x04
		req.DST_ADDR = []byte(ip6)
	} else {
		return req, errors.New("host error")
	}
	return req, nil
}

func (c *Conn) Close() error {
	return c.rw.Close()
}

func (c *Conn) Dial(network, address string) (net.Conn, error) {
	req, err := c.requests(network, address)
	if err != nil {
		return nil, err
	}

	log.Println("send:", req.ToByte())

	c.rw.Write(req.ToByte())
	buf := make([]byte, 231)
	n, err := c.rw.Read(buf)
	if err != nil {
		return nil, err
	}

	log.Println("recv:", buf[:n])

	rep := Replies{}
	rep.Serialization(buf[:n])
	if rep.REP == 0x00 {
		return c.rw, nil
	} else {
		return nil, fmt.Errorf("error replies REP [%d]", rep.REP)
	}
}
