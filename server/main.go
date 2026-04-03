package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

// bufferedConn makes bufio.Reader-readable bytes also available to net/http.
// net/http will read from Conn directly, so we override Read to read from the same bufio.Reader.
type bufferedConn struct {
	net.Conn
	r *bufio.Reader
}

func (c *bufferedConn) Read(p []byte) (int, error) {
	return c.r.Read(p)
}

// singleConnListener is a net.Listener that yields exactly one connection.
// It allows us to call http.Server.Serve(conn) for an already-accepted net.Conn.
type singleConnListener struct {
	conn net.Conn
	used bool
}

func (l *singleConnListener) Accept() (net.Conn, error) {
	if l.used {
		return nil, errors.New("single conn served")
	}
	l.used = true
	if l.conn == nil {
		return nil, errors.New("no conn")
	}
	return l.conn, nil
}

func (l *singleConnListener) Close() error { return nil }

func (l *singleConnListener) Addr() net.Addr {
	if l.conn == nil {
		return &net.TCPAddr{}
	}
	return l.conn.LocalAddr()
}

func handleTCPFirstByte0x05(conn net.Conn, br *bufio.Reader) {
	// Placeholder behavior:
	// - If you want to implement SOCKS5 server logic, replace this function.
	// - Here we just read the incoming bytes and log chunk sizes.
	log.Printf("[tcp] 0x05 branch: %s", conn.RemoteAddr())

	buf := make([]byte, 32*1024)
	for {
		n, err := br.Read(buf)
		if n > 0 {
			log.Printf("[tcp] read %d bytes", n)
		}
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Printf("[tcp] read error: %v", err)
			}
			return
		}
	}
}

func handleHTTP(conn net.Conn, br *bufio.Reader, httpSrv *http.Server) {
	// Clear the initial peek deadline; net/http will manage its own deadlines.
	_ = conn.SetReadDeadline(time.Time{})

	bc := &bufferedConn{Conn: conn, r: br}
	l := &singleConnListener{conn: bc}

	// Serve returns when the handler finishes and the listener refuses new conns.
	if err := httpSrv.Serve(l); err != nil {
		// singleConnListener returns a non-nil error on purpose to stop Serve.
		log.Printf("[http] serve end: %v", err)
	}
}

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	readTimeout := flag.Duration("read-timeout", 30*time.Second, "initial peek/read timeout")

	flag.Parse()

	// Example HTTP handler (uses http.DefaultServeMux so you can replace/extend with http.HandleFunc).
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = fmt.Fprintf(w, "ok method=%s path=%s\n", r.Method, r.URL.Path)
	})
	httpSrv := &http.Server{
		Addr:              *addr,
		Handler:           http.DefaultServeMux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("listen %s failed: %v", *addr, err)
	}
	defer ln.Close()
	log.Printf("listening on %s", *addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}

		go func(c net.Conn) {
			defer c.Close()

			_ = c.SetReadDeadline(time.Now().Add(*readTimeout))
			br := bufio.NewReader(c)

			first, err := br.Peek(1)
			if err != nil {
				log.Printf("[mux] peek error: %v (%s)", err, c.RemoteAddr())
				return
			}

			if first[0] == 0x05 {
				_ = c.SetReadDeadline(time.Now().Add(*readTimeout))
				handleTCPFirstByte0x05(c, br)
				return
			}

			handleHTTP(c, br, httpSrv)
		}(conn)
	}
}
