// Copyright (c) 2022 Bruce Zhang
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package evnet

import (
	"log"
	"math/rand"
	"net"
	"os"
	"runtime"

	. "github.com/bruce2233/evnet/socket"

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

type PollAttachment struct {
	Fd    int
	Event unix.EpollEvent
}

var subReactorBufferCap = 64 * 1024
var SubReactorsNum = 5

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

	log.Println("debug: ", "mainReactor create", "poller: ", mainReactor.poller, "linstenerFd is ", mainReactor.listener.Fd)

	for {
		log.Println("debug: ", "main reactor start polling")
		eventsList := mainReactor.poller.Polling()
		//working
		for _, v := range eventsList {
			nfd, raddr, err := AcceptSocket(int(v.Fd))
			log.Println("debug: ", nfd, raddr.Network(), raddr.String())
			os.NewSyscallError("AcceptSocket Err", err)
			//convert from nfd, tcpAddr to a socket
			// err = mainReactor.subReactors[0].poller.AddPollRead(int(v.Fd))
			idx := rand.Intn(2)
			laddr := mainReactor.listener.addr
			conn, _ := NewConn(nfd, laddr, raddr)
			log.Println("debug: ", "mainReactor accept connection ", raddr, "to subReactor ", idx)
			registerConn(mainReactor.subReactors[idx], conn)
		}
	}
}

func (sr *SubReactor) Loop() {
	for {
		log.Println("debug: ", "SubReactor start polling", sr.poller.Fd)
		eventList := sr.poller.Polling()
		for _, event := range eventList {
			log.Println("debug: ", "events: ", event.Events)
			if event.Events&OutEvents != 0 {
				c := sr.connections[int(event.Fd)]
				log.Println("debug: ", "subReactor OutEvents from: ", c.RemoteAddr())
				sr.write(c)
			}

			if event.Events&InEvents != 0 {
				// closeConn(sr,
				// sr.read(sr.connections[(int)(event.Fd)])
				c := sr.connections[int(event.Fd)]
				log.Println("debug: ", "subReactor InEvents from: ", c.RemoteAddr())
				sr.read(c)
			}
		}
	}
}

func (sr *SubReactor) Polling(callback func(c Conn)error){


}

func (sr *SubReactor) read(c *conn) error {
	n, err := unix.Read(c.Fd(), sr.buffer)
	log.Println("debug: ", "sr read n:", n)
	// log.Println("sr.buffer len():", len(sr.buffer))
	if n == 0 {
		sr.closeConn(c)
	}
	if err != nil {
		if err == unix.EAGAIN {
			log.Println("err: ", unix.EAGAIN)
			return unix.EAGAIN
		}
		return unix.ECONNRESET
	}
	c.inboundBuffer = sr.buffer[:n]
	(**sr.eventHandlerPP).OnTraffic(c)
	return nil
}

func (sr *SubReactor) write(c *conn) error {
	buffedLen := len(c.outboundBuffer)
	n, err := unix.Write(c.Fd(), c.outboundBuffer)
	if err != nil {
		log.Println("error: ", "subReactor Write error")
	}
	if n == buffedLen {
		c.outboundBuffer = c.outboundBuffer[n:]
		c.
	}
	c.outboundBuffer = c.outboundBuffer[n:]
	return err
}

func (poller *Poller) AddPollRead(pafd int) error {
	// log.Println("add poll read", "pollerFd: ", poller.Fd, "epolledFd: ", pafd)
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
	unix.Close(c.Fd())
	sr.poller.Delete(c.Fd())
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
	// log.Println("poller start waiting", poller.Fd)
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
