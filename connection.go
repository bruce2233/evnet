package evnet

import (
	"errors"
	"io"
	"net"

	. "github.com/bruce2233/evnet/socket"
	"golang.org/x/sys/unix"
)

type Conn interface {
	//working
	io.Reader
	Writer
	Socket

	Next(n int) ([]byte, error)

	// LocalAddr is the connection's local socket address.
	LocalAddr() (addr net.Addr)

	// RemoteAddr is the connection's remote peer address.
	RemoteAddr() (addr net.Addr)
}

type Writer interface {
	io.Writer

	AsyncWrite([]byte, func(Conn) error) error
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

//sync Write
func (c *conn) Write(p []byte) (n int, err error) {

	if len(c.outboundBuffer) > 0 {
		return -1, errors.New("Previous data waiting")
	}

	bufferedLen := len(c.outboundBuffer)

	sent, err := unix.Write(c.Fd(), p)
	for sent < len(c.outboundBuffer) {
		if sent < bufferedLen {
			c.outboundBuffer = c.outboundBuffer[sent:]
		}
	}
	return -1, err
}

//Async Write
func (c *conn) AsyncWrite(p []byte, AfterWritten func(c Conn) (err error)) error {

	if len(c.outboundBuffer) > 0 {
		return errors.New("Previous data waiting")
	}

	bufferedLen := len(c.outboundBuffer)

	sent, err := unix.Write(c.Fd(), p)
	for sent < len(c.outboundBuffer) {
		if sent < bufferedLen {
			c.outboundBuffer = c.outboundBuffer[sent:]
		}
	}
	return err
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
