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
	reactor        *SubReactor
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

	// if len(c.outboundBuffer) > 0 {
	// return -1, errors.New("Previous data waiting")
	// }
	sentSum := 0
	// bufferedLen := len(c.outboundBuffer)
	sent, err := unix.Write(c.Fd(), p)
	if sent != -1 {
		p = p[sent:]
		sentSum += sent
	}
	for sent > 0 {
		sent, err = unix.Write(c.Fd(), p)
		p = p[sent:]
		sentSum += sent
		// if sent < bufferedLen {
		// c.outboundBuffer = c.outboundBuffer[sent:]
		// }
	}
	return sentSum, err
}

//Async Write
func (c *conn) AsyncWrite(p []byte, AfterWritten func(c Conn) (err error)) error {

	if len(c.outboundBuffer) > 0 {
		return errors.New("Previous data waiting")
	}
	c.outboundBuffer = p
	c.reactor.poller.ModReadWrite(c.Fd())
	return nil
}

func (c *conn) Next(n int) (buf []byte, err error) {
	//working
	if n <= 0 {
		return c.inboundBuffer, nil
	}
	defer c.Discard(-1)
	return nil, nil
}

func (c *conn) Discard(n int)(int, error){
	if n == -1{
		discardedLen := len(c.inboundBuffer)
		c.inboundBuffer = make([]byte, 0)
		return discardedLen, nil
	}
	return 0,nil
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
