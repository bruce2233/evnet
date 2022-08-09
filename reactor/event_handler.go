package reactor

import (
)

type EventHandler interface {
	//working
	HandleConn(c Conn)
	GetHandleConn() func(Conn)
}

type BuiltinEventHandler struct {
}

func (beh BuiltinEventHandler) HandleConn(c Conn) {
	println("Handle Conn triggered ", c.Fd())
}

func (beh BuiltinEventHandler) GetHandleConn() func(Conn) {
	return beh.HandleConn
}
