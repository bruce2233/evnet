package reactor

import (
	. "evnet/connection"
)

type EventHandler interface {
	//working
	HandleConn(c Conn)
}

type BuiltinEventHandler struct {
}

func (beh BuiltinEventHandler) HandleConn(c Conn) {
	println("Handle Conn triggered ", c.Fd())
}
