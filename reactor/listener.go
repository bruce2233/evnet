package reactor

import "net"

type Listener struct {
	//working
	addr net.Addr
	Fd   int
}
