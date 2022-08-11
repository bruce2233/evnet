package evnet

import (
	"os"

	"golang.org/x/sys/unix"
)

type Poller struct {
	Fd int
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

// Delete removes the given file-descriptor from the poller.
func (p *Poller) Delete(fd int) error {
	return os.NewSyscallError("epoll_ctl del", unix.EpollCtl(p.Fd, unix.EPOLL_CTL_DEL, fd, nil))
}
