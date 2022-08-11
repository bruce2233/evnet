package evnet

import (
	"log"
	"math/rand"
	"net"
	"testing"
	"time"

	. "github.com/bruce2233/evnet/socket"

	"golang.org/x/sys/unix"
)

const streamLen = 1024 * 1024

type testServer struct {
	tester  *testing.T
	network string
	address string
}

type MyHandler struct {
	BuiltinEventHandler
}

func (mh MyHandler) OnTraffic(c Conn) {
	p, _ := c.Next(-1)
	log.Println("receive", len(p))
}

func (mh MyHandler) OnClose(c Conn) {
	log.Println("On Close Trigger")
}
func TestMainRec(t *testing.T) {
	ts := testServer{
		network: "tcp",
		address: "127.0.0.1:9000",
	}
	readHandler := MyHandler{BuiltinEventHandler{}}
	// Run(readHandler, ts.network+"://"+ts.address, WithLogPath("tmp.log"))
	Run(readHandler, ts.network+"://"+ts.address)
}

func TestClientWrite(t *testing.T) {
	startClient(t, "tcp", "127.0.0.1:9000")
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

func startClient(t *testing.T, network, address string) {
	reqData := make([]byte, streamLen)
	rand.Read(reqData)
	c, _ := net.Dial(network, address)
	defer c.Close()
	n, err := c.Write(reqData)
	t.Log(n, err)
	time.Sleep(1000)
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
