package socks5

import (
	"errors"
	"fmt"
	"net"
	"strconv"
)

// UDP Read Write func
type UDPClient struct {
	net.Conn
	DST_ADDR string
	DST_PORT uint16
}

func (c *UDPClient) Read(b []byte) (int, error) {
	buf := make([]byte, 1024)
	n, err := c.Conn.Read(buf)
	datagram := UDPDatagram{}
	if err != nil {
		return 0, err
	}
	err = datagram.Decode(buf[:n])
	if err != nil {
		return 0, err
	}
	bufLen := len(datagram.DATA)
	if len(b) < bufLen {
		return 0, fmt.Errorf("UDPClient Read buf si small[%d]", bufLen)
	}
	copy(b[:bufLen], datagram.DATA[:])
	return bufLen, err
}

func (c *UDPClient) Write(data []byte) (int, error) {
	datagram := UDPDatagram{
		DST_ADDR: c.DST_ADDR,
		DST_PORT: c.DST_PORT,
		DATA:     data,
	}
	if d, err := datagram.Encode(); err == nil {
		n, err := c.Conn.Write(d)
		if err != nil {
			return 0, err
		}
		if n != len(d) {
			return 0, errors.New("write byte len error")
		}
		return len(data), nil
	} else {
		return 0, err
	}
}

func (c *UDPClient) Close() error {
	return c.Conn.Close()
}

// Socks5 Clint Conn
type Conn struct {
	rw    net.Conn
	rwUDP *UDPClient
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
	case METHOD_USER:
		buf := []byte{0x01, byte(len(user))}
		buf = append(buf, []byte(user)...)
		buf = append(buf, byte(len(pass)))
		buf = append(buf, []byte(pass)...)
		if _, err := conn.Write(buf); err != nil {
			conn.Close()
			return Conn{}, err
		}
		resBuf := make([]byte, 2)
		rn, err := conn.Read(resBuf)
		if err != nil {
			conn.Close()
			return Conn{}, err
		}
		if rn != 2 || resBuf[0] != 0x01 {
			conn.Close()
			return Conn{}, errors.New("error socks5 username/password Version")
		}
		if resBuf[1] != 0x00 {
			conn.Close()
			return Conn{}, errors.New("error socks5 username/password")
		}
		return Conn{rw: conn}, err
	default:
		conn.Close()
		return Conn{}, fmt.Errorf("error socks method not match [%d]", buf[1])
	}
}

// Parse SOCKS5 Requests struct
// https://datatracker.ietf.org/doc/html/rfc1928#section-4
func (c *Conn) requests(network, address string) (Requests, error) {
	req := Requests{}

	if network == "tcp" {
		req.CMD = CMD_CONNECT
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
		req.DST_ADDR = host
	} else if network == "udp" {
		req.CMD = CMD_UDP_ASSOCIATE
		req.DST_PORT = 0
		req.DST_ADDR = "0.0.0.0"
	} else {
		return req, errors.New("not support " + network)
	}

	return req, nil
}

func (c *Conn) Close() error {
	if c.rwUDP != nil {
		c.rwUDP.Close()
	}
	return c.rw.Close()
}

func (c *Conn) Dial(network, address string) (net.Conn, error) {
	req, err := c.requests(network, address)
	if err != nil {
		return nil, err
	}
	if b, err := req.Encode(); err == nil {
		c.rw.Write(b)
	} else {
		return nil, err
	}
	buf := make([]byte, 231)
	n, err := c.rw.Read(buf)
	if err != nil {
		return nil, err
	}

	rep := Replies{}
	rep.Decode(buf[:n])
	if rep.REP == 0x00 {
		if req.CMD == CMD_UDP_ASSOCIATE {
			UDPrw, err := net.Dial("udp", fmt.Sprintf("%s:%d", rep.DST_ADDR, rep.DST_PORT))
			if err != nil {
				return nil, err
			}
			host, strPort, err := net.SplitHostPort(address)
			if err != nil {
				return nil, err
			}
			intPort, err := strconv.Atoi(strPort)
			if err != nil {
				return nil, err
			} else if intPort < 0 || intPort > 65535 {
				return nil, errors.New("port error [" + strPort + "]")
			}

			c.rwUDP = &UDPClient{
				Conn:     UDPrw,
				DST_ADDR: host,
				DST_PORT: uint16(intPort),
			}
			return c.rwUDP, err
		} else {
			// how about the bind?
			return c.rw, nil
		}
	} else {
		return nil, fmt.Errorf("error replies REP [%d]", rep.REP)
	}
}
