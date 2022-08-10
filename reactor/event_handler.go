package reactor

type EventHandler interface {
	//working
	OnConn(c Conn)
	OnClose(c Conn)
}

type BuiltinEventHandler struct {
}

func (builtinEventHandler BuiltinEventHandler) OnConn(c Conn) {
	println("OnConn triggered ", c.Fd())
}

func (builtinEventHandler BuiltinEventHandler) OnClose(c Conn) {
	println("OnClose triggered ", c.Fd())
}
