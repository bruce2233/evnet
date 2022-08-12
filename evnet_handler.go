package evnet

import "log"

type EventHandler interface {
	//working
	OnConn(c Conn)
	OnClose(c Conn)
	OnTraffic(c Conn)
}

type BuiltinEventHandler struct {
}

func (builtinEventHandler BuiltinEventHandler) OnConn(c Conn) {
	log.Println("OnConn triggered ", c.Fd())
}

func (builtinEventHandler BuiltinEventHandler) OnClose(c Conn) {
	log.Println("OnClose triggered ", c.Fd())
}

func (builtinEventHandler BuiltinEventHandler) OnTraffic(c Conn) {
	log.Println("OnClose triggered ", c.Fd())
}
