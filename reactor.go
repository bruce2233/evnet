// Copyright (c) 2022 Bruce Zhang
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package evnet

import (
	. "evnet/socket"
	"log"
	"math/rand"
	"net"
	"os"
	"runtime"

	"golang.org/x/sys/unix"
)

type eventLoop struct {
	eventHandler EventHandler
}

type MainReactor struct {
	subReactors    []*SubReactor
	poller         *Poller
	eventHandlerPP **EventHandler
	listener       Listener
}

type SubReactor struct {
	poller         *Poller
	eventHandlerPP **EventHandler
	connections    map[int]*conn
	buffer         []byte
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
	ErrEvents       = unix.EPOLLERR | unix.EPOLLHUP | unix.EPOLLRDHUP
	// OutEvents combines EPOLLOUT event and some exceptional events.
	OutEvents = ErrEvents | unix.EPOLLOUT
	// InEvents combines EPOLLIN/EPOLLPRI events and some exceptional events.
	InEvents = ErrEvents | unix.EPOLLIN | unix.EPOLLPRI
)

var subReactorBufferCap = 64 * 1024
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
	log.Println("init main poller: ", p)

	//init subReactor
	for i := 0; i < SubReactorsNum; i++ {
		newSubReactor := &SubReactor{}
		newSubReactor.poller, err = OpenPoller()
		newSubReactor.connections = make(map[int]*conn)
		newSubReactor.buffer = make([]byte, subReactorBufferCap)
		newSubReactor.eventHandlerPP = mr.eventHandlerPP
		mr.subReactors = append(mr.subReactors, newSubReactor)
	}
	if err != nil {
		return err
	}

	//init listen socket
	listenerFd, _, err := TcpSocket(proto, addr, true)
	mr.listener.Fd = listenerFd
	if err != nil {
		log.Println(err)
	}

	//add listenerFd epoll
	mr.poller.AddPollRead(mr.listener.Fd)

	//start subReactor
	for i := range mr.subReactors {
		go startSubReactor(mr.subReactors[i])
	}

	//setEventHandler
	beh := new(BuiltinEventHandler)
	setEventHandler(mr, beh)
	for i := range mr.subReactors {
		mr.subReactors[i].eventHandlerPP = mr.eventHandlerPP
	}
	return nil
}

func startSubReactor(sr *SubReactor) {
	// sr.poller.Polling()
	sr.Loop()
}

func (mr *MainReactor) SetEventHandler(eh EventHandler) {
	setEventHandler(mr, eh)
}

func setEventHandler(mr *MainReactor, eh EventHandler) {
	// mr.eventHandler.HandleConn = eh.GetHandleConn()
	if mr.eventHandlerPP == nil {
		mr.eventHandlerPP = new(*EventHandler)
		*mr.eventHandlerPP = &eh
	} else {
		*mr.eventHandlerPP = &eh
	}
}

func (mainReactor *MainReactor) Loop() {

	log.Println("poller: ", mainReactor.poller, "add poll linstenerFd", mainReactor.listener.Fd)

	for {
		log.Println("main reactor start loop")
		eventsList := mainReactor.poller.Polling()
		//working
		for _, v := range eventsList {
			nfd, raddr, err := AcceptSocket(int(v.Fd))
			log.Println(nfd)
			log.Println(raddr.Network())
			log.Println(raddr.String())
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

func (sr *SubReactor) Loop() {
	log.Println("SubReactor start polling", sr.poller.Fd)
	eventList := sr.poller.Polling()
	for _, event := range eventList {
		if event.Events&InEvents != 0 {
			// closeConn(sr,
			sr.read(sr.connections[(int)(event.Fd)])
		}
		log.Println(event.Fd, " event")
		conn := sr.connections[int(event.Fd)]
		sr.read(conn)
	}
}

func (sr *SubReactor) read(c *conn) error {
	n, err := unix.Read(c.Fd(), sr.buffer)
	log.Println("sr read n:", n)
	log.Println("sr.buffer len():", len(sr.buffer))
	if n == 0 {
		sr.closeConn(c)
	}
	if err != nil {
		if err == unix.EAGAIN {
			// log.Println(unix.EAGAIN)
			return unix.EAGAIN
		}
		return unix.ECONNRESET
	}
	c.inboundBuffer = sr.buffer[:n]
	(**sr.eventHandlerPP).OnConn(c)
	return nil
}

func (poller *Poller) AddPollRead(pafd int) error {
	log.Println("add poll read", "pollerFd: ", poller.Fd, "epolledFd: ", pafd)
	err := addPollRead(poller.Fd, pafd)
	return err
}

func addPollRead(epfd int, fd int) error {
	err := unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Fd: int32(fd), Events: readEvents})
	return os.NewSyscallError("AddPollRead Error", err)
}

func registerConn(sr *SubReactor, connItf Conn) error {
	// conn, err := NewConn(socket.Fd(),socket.)
	// subReactor[]
	sr.connections[connItf.Fd()] = connItf.(*conn)
	sr.poller.AddPollRead(connItf.Fd())
	return nil
}

func (sr *SubReactor) closeConn(c Conn) {
	//working
	(**sr.eventHandlerPP).OnClose(c)
	delete(sr.connections, c.Fd())
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
			log.Println("sa4 asset error")
		}
		return nfd, &net.TCPAddr{IP: sa4.Addr[:], Port: sa4.Port}, err
	}
	return -1, nil, err
}

//block
func (poller *Poller) Polling() []unix.EpollEvent {
	log.Println("poller start waiting", poller.Fd)
	eventList := polling(poller.Fd)
	return eventList
}

func polling(epfd int) []unix.EpollEvent {
	eventsList := make([]unix.EpollEvent, 128)

	var n int
	var err error
	for {
		n, err = unix.EpollWait(epfd, eventsList, -1)
		log.Println("epoll trigger")
		// if return EINTR, EpollWait again
		// debugging will trigger unix.EINTR error
		if n < 0 && err == unix.EINTR {
			runtime.Gosched()
			continue
		} else if err != nil {
			// logging.Errorf("error occurs in polling: %v", os.NewSyscallError("epoll_wait", err))
			panic("EpollWait Error")
		}
		break
	}
	return eventsList[:n]
}
