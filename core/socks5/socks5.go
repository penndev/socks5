// https://datatracker.ietf.org/doc/html/rfc1928
// https://datatracker.ietf.org/doc/html/rfc1929
package socks5

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
)

const Version = 0x05

// o  X'00' NO AUTHENTICATION REQUIRED
// o  X'01' GSSAPI
// o  X'02' USERNAME/PASSWORD
// o  X'03' to X'7F' IANA ASSIGNED
// o  X'80' to X'FE' RESERVED FOR PRIVATE METHODS
// o  X'FF' NO ACCEPTABLE METHODS
const (
	METHOD_NO_AUTH = 0x00
	METHOD_USER    = 0x02

	CMD_CONNECT       = 0x01
	CMD_UDP_ASSOCIATE = 0x03
)

type Requests struct {
	// o  CONNECT X'01'
	// o  BIND X'02'
	// o  UDP ASSOCIATE X'03'
	CMD byte
	//X'01'
	// the address is a Version-4 IP address, with a length of 4 octets
	//X'03'
	// the address field contains a fully-qualified domain name.  The first
	// octet of the address field contains the number of octets of name that
	// follow, there is no terminating NUL octet.
	//X'04'
	//the address is a Version-6 IP address, with a length of 16 octets.
	ATYP     byte
	DST_ADDR string
	DST_PORT uint16
}

func (r *Requests) Encode() ([]byte, error) {
	// parse ip
	var bufAddr []byte
	parseIP := net.ParseIP(r.DST_ADDR)
	if parseIP == nil {
		r.ATYP = 0x03
		bufAddr = []byte(r.DST_ADDR)
	} else if ip4 := parseIP.To4(); ip4 != nil {
		r.ATYP = 0x01
		bufAddr = []byte(ip4)
	} else if ip6 := parseIP.To16(); ip6 != nil {
		r.ATYP = 0x04
		bufAddr = []byte(ip6)
	} else {
		return []byte{}, errors.New("host error")
	}
	//  parse port
	bufPort := make([]byte, 2)
	binary.BigEndian.PutUint16(bufPort, uint16(r.DST_PORT))

	buf := []byte{Version, r.CMD, 0x00, r.ATYP}
	if r.ATYP == 0x03 {
		buf = append(buf, byte(len(r.DST_ADDR)))
	}

	buf = append(buf, bufAddr...)
	buf = append(buf, bufPort...)
	return buf, nil
}

type Replies struct {
	// o  X'00' succeeded
	// o  X'01' general SOCKS server failure
	// o  X'02' connection not allowed by ruleset
	// o  X'03' Network unreachable
	// o  X'04' Host unreachable
	// o  X'05' Connection refused
	// o  X'06' TTL expired
	// o  X'07' Command not supported
	// o  X'08' Address type not supported
	// o  X'09' to X'FF' unassigned
	REP byte
	//X'01'
	// the address is a Version-4 IP address, with a length of 4 octets
	//X'03'
	// the address field contains a fully-qualified domain name.  The first
	// octet of the address field contains the number of octets of name that
	// follow, there is no terminating NUL octet.
	//X'04'
	//the address is a Version-6 IP address, with a length of 16 octets.
	ATYP     byte
	DST_ADDR string
	DST_PORT uint16
}

func (r *Replies) Decode(buf []byte) error {
	if len(buf) < 7 {
		return errors.New("Replies Decode error byte")
	}
	if buf[0] != Version {
		return fmt.Errorf("error version %d", buf[0])
	}
	r.REP = buf[1]
	r.ATYP = buf[3]
	portLen := 0
	if r.ATYP == 0x01 {
		if len(buf) != 10 {
			return errors.New("Replies Decode 0x01 err")
		}
		r.DST_ADDR = net.IP(buf[4:8]).String()
		portLen = 8
	} else if r.ATYP == 0x03 {
		addrLen := int(buf[4])
		if len(buf) != (addrLen + 7) {
			return errors.New("Replies Decode 0x03 err")
		}
		r.DST_ADDR = net.IP(buf[5 : 5+addrLen]).String()
		portLen = 5 + addrLen
	} else if r.ATYP == 0x04 {
		if len(buf) != 22 {
			return errors.New("Replies Decode 0x04 err")
		}
		r.DST_ADDR = net.IP(buf[4:21]).String()
		portLen = 21
	} else {
		return fmt.Errorf("Replies Decode atyp error %d", r.ATYP)
	}
	r.DST_PORT = binary.BigEndian.Uint16([]byte{buf[portLen], buf[portLen+1]})
	return nil
}

type UDPDatagram struct {
	// RSV | FRAG |
	//X'01'
	// the address is a Version-4 IP address, with a length of 4 octets
	//X'03'
	// the address field contains a fully-qualified domain name.  The first
	// octet of the address field contains the number of octets of name that
	// follow, there is no terminating NUL octet.
	//X'04'
	//the address is a Version-6 IP address, with a length of 16 octets.
	ATYP     byte
	DST_ADDR string
	DST_PORT uint16
	DATA     []byte
}

func (r *UDPDatagram) Decode(buf []byte) error {
	if len(buf) < 7 {
		return errors.New("UDPDatagram Decode error buf size")
	}
	r.ATYP = buf[3]
	portLen := 0
	if r.ATYP == 0x01 {
		if len(buf) < 10 {
			return errors.New("UDPDatagram Decode 0x01 err")
		}
		r.DST_ADDR = net.IP(buf[4:8]).String()
		portLen = 8
	} else if r.ATYP == 0x03 {
		addrLen := int(buf[4])
		if len(buf) < (addrLen + 7) {
			return errors.New("UDPDatagram Decode 0x03 err")
		}
		r.DST_ADDR = net.IP(buf[5 : 5+addrLen]).String()
		portLen = 5 + addrLen
	} else if r.ATYP == 0x04 {
		if len(buf) < 22 {
			return errors.New("UDPDatagram Decode 0x04 err")
		}
		r.DST_ADDR = net.IP(buf[4:21]).String()
		portLen = 21
	} else {
		log.Println(buf)
		return fmt.Errorf("UDPDatagram Decode atyp error %d", r.ATYP)
	}
	r.DST_PORT = binary.BigEndian.Uint16([]byte{buf[portLen], buf[portLen+1]})
	r.DATA = buf[portLen+2:]
	return nil
}

func (r *UDPDatagram) Encode() ([]byte, error) {
	// parse ip
	var bufAddr []byte
	parseIP := net.ParseIP(r.DST_ADDR)
	if parseIP == nil {
		r.ATYP = 0x03
		bufAddr = []byte(r.DST_ADDR)
	} else if ip4 := parseIP.To4(); ip4 != nil {
		r.ATYP = 0x01
		bufAddr = []byte(ip4)
	} else if ip6 := parseIP.To16(); ip6 != nil {
		r.ATYP = 0x04
		bufAddr = []byte(ip6)
	} else {
		return []byte{}, errors.New("host error")
	}
	//  parse port
	bufPort := make([]byte, 2)
	binary.BigEndian.PutUint16(bufPort, r.DST_PORT)

	buf := []byte{0x00, 0x00, 0x00, r.ATYP}
	if r.ATYP == 0x03 {
		buf = append(buf, byte(len(r.DST_ADDR)))
	}
	buf = append(buf, bufAddr...)
	buf = append(buf, bufPort...)
	buf = append(buf, r.DATA...)
	return buf, nil
}
