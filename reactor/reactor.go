// MIT License

// Copyright (c) [2022] [Bruce Zhang]

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package reactor

import (
	. "evnet/connection"
	. "evnet/socket"
	"math/rand"
	"net"
	"os"
	"runtime"

	"golang.org/x/sys/unix"
)

type MainReactor struct {
	subReactors  []*SubReactor
	poller       *Poller
	eventHandler EventHandler
	listener     Listener
}

type SubReactor struct {
	poller       *Poller
	eventHandler EventHandler
	connections  map[int]Conn
}

type Poller struct {
	Fd  int
	efd int
}
type PollAttachment struct {
	Fd    int
	Event unix.EpollEvent
}

const (
	readEvents      = unix.EPOLLPRI | unix.EPOLLIN
	writeEvents     = unix.EPOLLOUT
	readWriteEvents = readEvents | writeEvents
)

var SubReactorsNum = 2

//create a new poller
func OpenPoller() (poller *Poller, err error) {
	poller = new(Poller)
	if poller.Fd, err = unix.EpollCreate1(unix.EPOLL_CLOEXEC); err != nil {
		poller = nil
		err = os.NewSyscallError("epoll_create1", err)
		return
	}
	// if poller.efd, err = unix.Eventfd(0, unix.EFD_NONBLOCK|unix.EFD_CLOEXEC); err != nil {
	// 	// _ = poller.Close()
	// 	poller = nil
	// 	err = os.NewSyscallError("eventfd", err)
	// 	return
	// }
	// poller.efdBuf = make([]byte, 8)
	// if err = poller.AddRead(&PollAttachment{FD: poller.efd}); err != nil {
	// 	_ = poller.Close()
	// 	poller = nil
	// 	return
	// }
	// poller.asyncTaskQueue = queue.NewLockFreeQueue()
	// poller.urgentAsyncTaskQueue = queue.NewLockFreeQueue()
	return
}

//The proto paramater MUST be "tcp"
//The addr parameter
func (mainReactor *MainReactor) Listen(proto, addr string) (fd int, netAddr net.Addr, err error) {
	fd, netAddr, err = TcpSocket(proto, addr, true)
	return
}

//Create poller and listener for mainReactor
func (mr *MainReactor) Init(proto, addr string) error {
	//init poller
	p, err := OpenPoller()
	if err != nil {
		return err
	}
	mr.poller = p
	println("init main poller: ", p)

	//init subReactor
	for i := 0; i < SubReactorsNum; i++ {
		newSubReactor := &SubReactor{}
		newSubReactor.poller, err = OpenPoller()
		newSubReactor.connections = make(map[int]Conn)
		mr.subReactors = append(mr.subReactors, newSubReactor)
	}
	if err != nil {
		return err
	}

	//init listen socket
	listenerFd, _, err := TcpSocket(proto, addr, true)
	mr.listener.Fd = listenerFd
	if err != nil {
		println(err)
	}

	//add listenerFd epoll
	mr.poller.AddPollRead(mr.listener.Fd)

	//start subReactor
	for i := range mr.subReactors {
		go startSubReactor(mr.subReactors[i])
	}

	//setEventHandler
	beh := new(BuiltinEventHandler)
	mr.eventHandler = beh
	for i := range mr.subReactors {
		mr.subReactors[i].eventHandler = beh
	}
	return nil
}

func startSubReactor(sr *SubReactor) {
	println("SubReactor start polling", sr.poller.Fd)
	eventList := sr.poller.Polling()
	for _, event := range eventList {
		println(event.Fd, " event")
		conn := sr.connections[int(event.Fd)]
		sr.eventHandler.HandleConn(conn)
	}
}

func (mr *MainReactor) SetEventHandler(eh EventHandler) {
	setEventHandler(mr, eh)
}

func setEventHandler(mr *MainReactor, eh EventHandler) {
	mr.eventHandler = eh
}

func (mainReactor *MainReactor) Loop() {

	println("poller: ", mainReactor.poller, "add poll linstenerFd", mainReactor.listener.Fd)

	for {
		println("main reactor start loop")
		eventsList := mainReactor.poller.Polling()
		//working
		for _, v := range eventsList {
			nfd, raddr, err := AcceptSocket(int(v.Fd))
			println(nfd)
			println(raddr.Network())
			println(raddr.String())
			os.NewSyscallError("AcceptSocket Err", err)
			//convert from nfd, tcpAddr to a socket
			// err = mainReactor.subReactors[0].poller.AddPollRead(int(v.Fd))
			idx := rand.Intn(2)
			laddr := mainReactor.listener.addr
			conn, _ := NewConn(nfd, laddr, raddr)
			registerConn(mainReactor.subReactors[idx], conn)
		}
	}
}

func (poller *Poller) AddPollRead(pafd int) error {
	println("add poll read", "pollerFd: ", poller.Fd, "epolledFd: ", pafd)
	err := addPollRead(poller.Fd, pafd)
	return err
}

func addPollRead(epfd int, fd int) error {
	err := unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Fd: int32(fd), Events: readEvents})
	return os.NewSyscallError("AddPollRead Error", err)
}

func registerConn(sr *SubReactor, conn Conn) error {
	// conn, err := NewConn(socket.Fd(),socket.)
	// subReactor[]
	sr.connections[conn.Fd()] = conn
	sr.poller.AddPollRead(conn.Fd())
	return nil
}

func AcceptSocket(fd int) (int, net.Addr, error) {
	nfd, sa, err := unix.Accept(fd)
	if err = os.NewSyscallError("fcntl nonblock", unix.SetNonblock(nfd, true)); err != nil {
		return -1, nil, err
	}
	switch sa.(type) {
	case *unix.SockaddrInet4:
		sa4, ok := sa.(*unix.SockaddrInet4)
		if !ok {
			println("sa4 asset error")
		}
		return nfd, &net.TCPAddr{IP: sa4.Addr[:], Port: sa4.Port}, err
	}
	return -1, nil, err
}

//block
func (poller *Poller) Polling() []unix.EpollEvent {
	eventsList := make([]unix.EpollEvent, 128)

	println("poller start waiting", poller.Fd)
	var n int
	var err error
	for {
		n, err = unix.EpollWait(poller.Fd, eventsList, -1)
		println("epoll trigger")
		// if return EINTR, EpollWait again
		// debugging will trigger unix.EINTR error
		if n < 0 && err == unix.EINTR {
			runtime.Gosched()
			continue
		}
		println("epoll events num: ", n, err)
		if err != nil {
			println(err)
		}
		break
	}
	return eventsList[:n]
}
