# ⚡ Introduction
![](https://s3.bmp.ovh/imgs/2022/08/09/7704dbff7cbff64c.png)

`Evnet` is event-driven net framework with high performance deriving from gnet. It supports  [epoll](https://en.wikipedia.org/wiki/EpollNetty ) syscall in Linux only now.

# ✨ Features

- [x] Supporting an event-driven mechanism: epoll on Linux.
- [ ] Supporting multiple protocols/IPC mechanism: TCP, UDP and Unix Domain Socket.
- [ ] Supporting multiple load-balancing algorithms: Round-Robin, Source-Addr-Hash and Least-Connections
- [ ]  Efficient, reusable and elastic memory buffer and zero copy.

# 🎬 Getting started
```powershell 
go get github.com/bruce2233/evnet
```

```go
type MyHandler struct {
	BuiltinEventHandler
}

func (mh MyHandler) HandleConn(c Conn) {
	p := make([]byte, 128)
	n, err := unix.Read(c.Fd(), p)
	println("unix.Read bytes n: ", n)
	if err != nil {
		panic("unix.Read")
	}
	println(string(p[:n]))
}

//main.go
    readHandler := MyHandler{BuiltinEventHandler{}}
	m := new(MainReactor)
	m.Init("tcp", "127.0.0.1:9000")
	m.SetEventHandler(readHandler)
	m.Loop()
```