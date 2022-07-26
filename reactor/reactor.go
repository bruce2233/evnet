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

	"golang.org/x/sys/unix"
)

type MainReactor struct {
	subReactors  []*SubReactor
	poller       *Poller
	listenerFd   int
	eventHandler EventHandler
}

type SubReactor struct {
	poller       *Poller
	eventHandler EventHandler
	connections  map[int]Conn
}

type EventHandler interface {
	//working
	HandleConn(conn Conn)
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
func (mainReactor *MainReactor) Init(proto, addr string) error {
	//init poller
	p, err := OpenPoller()
	if err != nil {
		return err
	}
	mainReactor.poller = p
	println("init main poller: ", p)

	//init subReactor
	for i := 0; i < SubReactorsNum; i++ {
		newSubReactor := &SubReactor{}
		newSubReactor.poller, err = OpenPoller()
		mainReactor.subReactors = append(mainReactor.subReactors, newSubReactor)
	}
	if err != nil {
		return err
	}
	println("init sub poller")

	// //start subReactor
	// for i := 0; i < SubReactorsNum; i++ {
	// 	go startSubReactor(mainReactor.subReactors[i])
	// }

	//init listen socket
	listenerFd, _, err := mainReactor.Listen(proto, addr)
	mainReactor.listenerFd = listenerFd
	println("set listener fd: ", listenerFd)
	if err != nil {
		return err
	}

	return nil
}

func startSubReactor(subReactor *SubReactor) {
	eventList := subReactor.poller.Polling()
	for _, event := range eventList {
		conn := subReactor.connections[int(event.Fd)]
		subReactor.eventHandler.HandleConn(conn)
	}
}

func (mainReactor *MainReactor) Loop() {
	//add listenerFd epoll
	mainReactor.poller.AddPollRead(mainReactor.listenerFd)
	println("poller: ", mainReactor.poller, "add poll linstenerFd", mainReactor.listenerFd)

	for {
		println("start polling")
		eventsList := mainReactor.poller.Polling()
		//working
		for _, v := range eventsList {
			nfd, tcpAddr, err := AcceptSocket(int(v.Fd))
			println(nfd)
			println(tcpAddr.Network())
			println(tcpAddr.String())
			os.NewSyscallError("AcceptSocket Err", err)
			//convert from nfd, tcpAddr to a socket
			// err = mainReactor.subReactors[0].poller.AddPollRead(int(v.Fd))
			idx := rand.Intn(2)

			addPollRead(mainReactor.subReactors[idx].poller.Fd, nfd)
		}
	}
}

func (poller *Poller) AddPollRead(pafd int) error {
	err := addPollRead(poller.Fd, pafd)
	return err
}
func addPollRead(epfd int, fd int) error {
	err := unix.EpollCtl(fd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Fd: int32(fd), Events: readEvents})
	return os.NewSyscallError("AddPollRead Error", err)
}

func registerConn(subReactor SubReactor, socket Socket) error {
	// conn, err := NewConn(socket.Fd(),socket.)
	// subReactor[]
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

	println(poller, "start waiting")
	n, err := unix.EpollWait(poller.Fd, eventsList, -1)
	println("epoll events num: ", n, err)
	if err != nil {
		println(err)
	}
	// for i := 0; i < n; i++ {
	// 	println(eventsList[i].Fd, eventsList[i].Events)
	// }
	return eventsList[:n]
}
