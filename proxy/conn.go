package proxy

import (
	"bufio"
	"net"
	"sync/atomic"
)

type Conn struct {
	net.Conn
	reader     *bufio.Reader
	readBytes  *uint64
	writeBytes *uint64
}

func (c *Conn) Read(p []byte) (n int, err error) {
	n, err = c.reader.Read(p)
	if n > 0 && c.readBytes != nil {
		atomic.AddUint64(c.readBytes, uint64(n))
	}
	return
}

func (c *Conn) Write(p []byte) (n int, err error) {
	n, err = c.Conn.Write(p)
	if n > 0 && c.writeBytes != nil {
		atomic.AddUint64(c.writeBytes, uint64(n))
	}
	return
}

func (c *Conn) Peek(n int) ([]byte, error) {
	return c.reader.Peek(n)
}

func NewConn(conn net.Conn, readBytes, writeBytes *uint64) *Conn {
	return &Conn{
		Conn:       conn,
		reader:     bufio.NewReader(conn),
		readBytes:  readBytes,
		writeBytes: writeBytes,
	}
}
