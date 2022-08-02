package test

import (
	. "evnet/reactor"
	// . "evnet/socket"
	"net"
	"testing"

	"golang.org/x/sys/unix"
)

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
}

func TestMainRec(t *testing.T) {
	m := new(MainReactor)
	m.Init("tcp", "127.0.0.1:9000")
	m.Loop()
}

func TestStartSubReactor(t *testing.T) {
	m := MainReactor{}
	m.Init()
}
func TestEpollListener(t *testing.T) {
	// socketFd, netAddr, err := TcpSocket("tcp", "127.0.0.1:8866", true)
	m := new(MainReactor)
	socketFd, netAddr, _ := m.Listen("tcp", "127.0.0.1:9000")
	t.Log(socketFd)
	t.Log(netAddr)
	poller, _ := OpenPoller()
	// epollFd, err := unix.EpollCreate1(unix.EPOLL_CLOEXEC)

	t.Log(poller.Fd)

	// unix.EpollCtl(poller.Fd, unix.EPOLL_CTL_ADD, socketFd, &unix.EpollEvent{
	// 	Fd: int32(socketFd), Events: unix.EPOLLIN})
	poller.AddPollRead(socketFd)
	el := make([]unix.EpollEvent, 128)
	t.Log("before epoll wait")
	n, err := unix.EpollWait(poller.Fd, el, -1)

	acceptFd := el[0].Fd
	ev := el[0].Events
	t.Log(acceptFd, ev)
	t.Log(err)
	t.Log(n)

	// bytes := make([]byte, 256)
	// bytesNum, err := unix.Read(int(acceptFd), bytes)
	// if err != nil {
	// 	t.Log(err)
	// }
	// t.Log(bytesNum)
	// t.Log(string(bytes))

	nfd, sa, err := unix.Accept(int(acceptFd))
	if err != nil {
		t.Log(err)
	}
	t.Log(nfd)
	t.Log(sa)
}

