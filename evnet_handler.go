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
	OnShutdown(mr *MainReactor) error
	OnBoot(mr *MainReactor) error
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

func (builtinEventHandler BuiltinEventHandler) OnShutdown(mr *MainReactor) error {
	log.Fatalln("OnShutdown triggered ")
	return nil
}

func (builtinEventHandler BuiltinEventHandler) OnBoot(mr *MainReactor) error {
	log.Fatalln("OnBoot triggered ")
	return nil
}
