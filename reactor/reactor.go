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
	"evnet/socket"
	"net"
	"os"

	"golang.org/x/sys/unix"
)

type MainReactor struct {
	subReactors  []SubReactor
	poller       *Poller
	listenerFd   int
	eventHandler EventHandler
}

type SubReactor struct {
	poller       *Poller
	eventHandler EventHandler
}

type EventHandler interface {
	//working
	HandleConn(conn *net.Conn)
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
	fd, netAddr, err = socket.TcpSocket(proto, addr, true)
	return
	net.Listen()
}

func (mainReactor *MainReactor) SetListenerFd(fd int) {
	mainReactor.listenerFd = fd
}
func (mainReactor *MainReactor) Init(proto, addr string) error {
	//init poller
	p, err := OpenPoller()
	if err != nil {
		return err
	}
	mainReactor.poller = p
	println("init poller poller: ", p)

	//init listen socket
	fd, _, err := mainReactor.Listen(proto, addr)
	mainReactor.SetListenerFd(fd)
	println("set listener fd: ", fd)
	if err != nil {
		return err
	}
	return nil
}

func (mainReactor *MainReactor) Loop() {
	mainReactor.poller.AddPollRead(mainReactor.listenerFd)
	println("poller: ", mainReactor.poller, "add poll linstenerFd", mainReactor.listenerFd)
	for {
		println("start polling")
		eventsList := mainReactor.poller.Polling()
		//working
		for _, v := range eventsList {
			tcpAddr, err := AcceptSocket(int(v.Fd))
			println(tcpAddr.Network())
			println(tcpAddr.String())
			os.NewSyscallError("AcceptSocket Err", err)
			// err := mainReactor.subReactors[0].poller.AddPollRead(int(v.Fd))
		}
	}
}

func (poller *Poller) AddPollRead(pafd int) error {
	err := unix.EpollCtl(poller.Fd, unix.EPOLL_CTL_ADD, pafd, &unix.EpollEvent{Fd: int32(pafd), Events: readEvents})
	return os.NewSyscallError("AddPollRead Error", err)
}

func AcceptSocket(fd int) (net.Addr, error) {
	nfd, sa, err := unix.Accept(fd)
	if err = os.NewSyscallError("fcntl nonblock", unix.SetNonblock(nfd, true)); err != nil {
		return nil, err
	}
	switch sa.(type) {
	case *unix.SockaddrInet4:
		sa4, ok := sa.(*unix.SockaddrInet4)
		if !ok {
			println("sa4 asset error")
		}
		return &net.TCPAddr{IP: sa4.Addr[:], Port: sa4.Port}, err
	}
	return nil, nil
}

func (poller *Poller) Polling() []unix.EpollEvent {
	eventsList := make([]unix.EpollEvent, 128)

	println(poller, "start waiting")
	n, err := unix.EpollWait(poller.Fd, eventsList, -1)
	println("epoll events num: ", n, err)
	// for i := 0; i < n; i++ {
	// 	println(eventsList[i].Fd, eventsList[i].Events)

	// }
	return eventsList[:n]
}
