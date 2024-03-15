// https://datatracker.ietf.org/doc/html/rfc1928
// https://datatracker.ietf.org/doc/html/rfc1929
package socks5

import (
	"encoding/binary"
	"fmt"
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
	DST_ADDR []byte
	DST_PORT uint16
}

func (r *Requests) ToByte() []byte {
	bufPort := make([]byte, 2)
	binary.BigEndian.PutUint16(bufPort, uint16(r.DST_PORT))
	buf := []byte{Version, r.CMD, 0x00, r.ATYP}
	if r.ATYP == 0x03 {
		buf = append(buf, byte(len(r.DST_ADDR)))
	}
	buf = append(buf, r.DST_ADDR...)
	buf = append(buf, bufPort...)
	return buf
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
	DST_ADDR []byte
	DST_PORT uint16
}

func (r *Replies) Serialization(buf []byte) error {
	if buf[0] != Version {
		return fmt.Errorf("error version %d", buf[0])
	}
	r.REP = buf[1]
	return nil
}
