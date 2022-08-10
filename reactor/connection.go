package reactor

import (
	. "evnet/socket"
	"io"
	"net"
)

type Conn interface {
	//working
	io.Reader
	io.Writer
	Socket

	Next(n int) ([]byte, error)

	// LocalAddr is the connection's local socket address.
	LocalAddr() (addr net.Addr)

	// RemoteAddr is the connection's remote peer address.
	RemoteAddr() (addr net.Addr)
}

type conn struct {
	fd             int
	localAddr      net.Addr
	remoteAddr     net.Addr
	inboundBuffer  []byte
	outboundBuffer []byte
}

func (c conn) Read(p []byte) (n int, err error) {
	//working
	if err != nil {
		return -1, err
	}
	return
}

func (c *conn) Write(p []byte) (n int, err error) {
	//working
	return -1, err
}

func (c *conn) Next(n int) (buf []byte, err error) {
	//working
	if n <= 0 {
		return c.inboundBuffer, nil
	}
	return nil, nil
}

func (c conn) Fd() int {
	return c.fd
}

func (c conn) LocalAddr() net.Addr {
	return c.localAddr
}

func (c conn) RemoteAddr() net.Addr {
	return c.remoteAddr
}

func NewConn(fd int, la net.Addr, ra net.Addr) (Conn, error) {
	conn := new(conn)
	conn.fd = fd
	conn.localAddr = la
	conn.remoteAddr = ra
	return conn, nil
}
