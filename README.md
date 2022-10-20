# ⚡ Introduction
![](https://s3.bmp.ovh/imgs/2022/08/09/7704dbff7cbff64c.png)

`Evnet` is event-driven net framework with high performance deriving from gnet and evio. It supports  [epoll](https://en.wikipedia.org/wiki/EpollNetty ) syscall in Linux only now. I'll complete docs continuously. If it's useful to you, Please STAR it!

# ✨ Features

- [x] Supporting an event-driven mechanism: epoll on Linux.
- [ ] Supporting multiple protocols/IPC mechanism: TCP, UDP and Unix Domain Socket.
- [ ]  Efficient, reusable and elastic memory buffer and zero copy.
- [ ] Supporting multiple load-balancing algorithms: Round-Robin, Source-Addr-Hash and Least-Connections

# 🎬 Getting started
```powershell 
go get github.com/bruce2233/evnet
```

```go
type MyHandler struct {
	BuiltinEventHandler
}

func (mh MyHandler) OnConn(c Conn) {
	p, _ := c.Next(-1)
	println(string(p))
}

func (mh MyHandler) OnClose(c Conn) {
	println("On Close Trigger")
}

//main.go
func TestMainRec(t *testing.T) {
	readHandler := MyHandler{BuiltinEventHandler{}}
	Run(readHandler, "tcp://127.0.0.1:9000")
}
```
