package evnet

import (
	"errors"
	"log"
	"math/rand"
	"net"
	"reflect"
	"testing"

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
	c.Write(p)
	log.Println("receive", len(p))
}

func (mh MyHandler) OnClose(c Conn) {
	log.Println("On Close Trigger")
}
func TestMainRec(t *testing.T) {
	ts := testServer{
		network: "tcp",
		address: "192.168.87.141:9000",
	}
	readHandler := MyHandler{BuiltinEventHandler{}}
	// Run(readHandler, ts.network+"://"+ts.address, WithLogPath("tmp.log"))
	Run(readHandler, ts.network+"://"+ts.address)
}

const (
	CLIENTNUM = 10
)

func TestClientWrite(t *testing.T) {
	done := make(chan bool)
	const reqDataSize int = 10240
	for i := 0; i < CLIENTNUM; i++ {
		go func() {
			reqData := make([]byte, reqDataSize)
			rand.Read(reqData)
			t.Log(reqData)
			c, err := net.Dial("tcp", "192.168.87.141:9000")
			if err != nil {
				t.Log(err)
			}
			c.Write(reqData)

			respData := make([]byte, reqDataSize)
			c.Read(respData)
			t.Log(respData)
			if !reflect.DeepEqual(reqData, respData) {
				t.Log(errors.New("write data and read data error"))
			}
			done <- true
		}()
	}
	// time.Sleep(10 * time.Second)
	ack := CLIENTNUM
	for ack > 0 {
		if <-done {
			ack--
			t.Log(ack, " client remains...")
		}
	}
}

func TestNet(t *testing.T) {
	conn, err := net.Dial("tcp", "192.168.87.141:9000")
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
	// total := 0
	// for {
	// 	inData := make([]byte, 10*1024)
	// 	conn.Read(inData)
	// 	t.Log("new: ", len(inData), "total: ", total)
	// 	total += len(inData)
	// }
	// conn.Close()
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
