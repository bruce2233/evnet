package reactor

import (
	"net"
	"testing"
)

func TestMainRec(t *testing.T) {
	m := new(MainReactor)
	m.Init("tcp", "127.0.0.1:9000")
	m.Loop()
}

func TestXxx(t *testing.T) {
	conn := net.Listen()
}
