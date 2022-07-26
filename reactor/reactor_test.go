package reactor

import "testing"

func TestMainRec(t *testing.T) {
	m := new(MainReactor)
	m.Init("tcp", "127.0.0.1:9000")
	m.Loop()
}
