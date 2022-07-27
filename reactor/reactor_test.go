package reactor

import (
	. "evnet/socket"
	"net"
	"testing"

	"golang.org/x/sys/unix"
)

func TestMainRec(t *testing.T) {
	m := new(MainReactor)
	m.Init("tcp", "127.0.0.1:9000")
	m.Loop()
}

func TestNet(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		t.Log(err)
	}
	t.Log(conn)
	bytes := []byte("Hello gnet")
	n, err := conn.Write(bytes)
	if err != nil {
		t.Log(err)
	}
	t.Log(n)
	conn.Close()
}

func TestLn(t *testing.T) {
	ln, _ := net.Listen("tcp", "127.0.0.1:9000")
	conn, _ := ln.Accept()
	p := make([]byte, 128)
	conn.Read(p)
	t.Log(string(p))
}

func TestEpollWait(t *testing.T) {
	epfd, _ := unix.EpollCreate1(unix.EPOLL_CLOEXEC)
	eventList := make([]unix.EpollEvent, 128)
	listenerFd, _, _ := TcpSocket("tcp", "127.0.0.1:9000", true)
	unix.SetNonblock(listenerFd, true)
	unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, listenerFd, &unix.EpollEvent{Fd: int32(listenerFd), Events: unix.EPOLLIN})
	// addPollRead(epfd, listenerFd)
	n, err := unix.EpollWait(epfd, eventList, -1)
	for i := 0; i < n; i++ {
		t.Log(eventList[i])
	}
	t.Log(n)
	if err != nil {
		t.Log(err)
	}
	t.Log("wait over")
	p := []byte{0, 1, 2}
	unix.Read(listenerFd, p)
	t.Log(string(p))
}
