package evnet

import (
	"errors"

	log "github.com/sirupsen/logrus"
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
	log.Info("OnConn triggered ", c.Fd())
	return nil
}

func (builtinEventHandler BuiltinEventHandler) OnClose(c Conn) error {
	log.Info("OnClose triggered ", c.Fd())
	return nil
}

func (builtinEventHandler BuiltinEventHandler) OnTraffic(c Conn) error {
	log.Debug("OnClose triggered ", c.Fd())
	return nil
}

func (builtinEventHandler BuiltinEventHandler) OnOpen(c Conn) error {
	log.Info("OnOpen triggered ", c.Fd())
	return nil
}

func (builtinEventHandler BuiltinEventHandler) OnShutdown(mr *MainReactor) error {
	log.Info("OnShutdown triggered ")
	return nil
}

func (builtinEventHandler BuiltinEventHandler) OnBoot(mr *MainReactor) error {
	log.Info("OnBoot triggered ")
	return nil
}
