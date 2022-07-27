package connection

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
}

type conn struct {
	fd         int
	localAddr  net.Addr
	remoteAddr net.Addr
}

func (conn *conn) Read(p []byte) (n int, err error) {
	//working
	return -1, err
}

func (conn *conn) Write(p []byte) (n int, err error) {
	//working
	return -1, err
}

func (conn *conn) Fd() int {
	return conn.fd
}

func NewConn(fd int, la net.Addr, ra net.Addr) (Conn, error) {
	conn := new(conn)
	conn.fd = fd
	conn.localAddr = la
	conn.remoteAddr = ra
	return conn, nil
}
