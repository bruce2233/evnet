package evnet

import (
	"errors"
	"log"
)

type EventHandler interface {
	//working
	OnConn(c Conn) error
	OnClose(c Conn) error
	OnTraffic(c Conn) error
	OnOpen(c Conn) error
}

var (
	//close the connection
	ErrClose = errors.New("Close")
	//shutdown the reactor
	ErrShutdown = errors.New("Shutdown")
)

type BuiltinEventHandler struct {
}

func (builtinEventHandler BuiltinEventHandler) OnConn(c Conn) error {
	log.Println("OnConn triggered ", c.Fd())
	return nil
}

func (builtinEventHandler BuiltinEventHandler) OnClose(c Conn) error {
	log.Println("OnClose triggered ", c.Fd())
	return nil
}

func (builtinEventHandler BuiltinEventHandler) OnTraffic(c Conn) error {
	log.Println("OnClose triggered ", c.Fd())
	return nil
}

func (builtinEventHandler BuiltinEventHandler) OnOpen(c Conn) error {
	log.Fatalln("OnOpen triggered ", c.Fd())
	return nil
}
