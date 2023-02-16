# ⚡ Introduction
![](https://s3.bmp.ovh/imgs/2022/08/09/7704dbff7cbff64c.png)

`Evnet` is event-driven net framework with high performance deriving from gnet and evio. It supports  [epoll](https://en.wikipedia.org/wiki/EpollNetty ) syscall in Linux. If it's useful to you, Please STAR it!

# ✨ Features

- [x] Supporting an event-driven mechanism: epoll on Linux.
- [ ] Supporting multiple protocols/IPC mechanism: TCP, UDP and Unix Domain Socket.
- [x] Efficient memory buffer and zero copy.
- [ ] Supporting multiple load-balancing algorithms: Round-Robin, Source-Addr-Hash and Least-Connections
- [ ] Implementation of gnet Client

# 🎬 Getting started
```powershell 
go get github.com/bruce2233/evnet
```

```go
import "github.com/bruce2233/evnet"

type MyHandler struct {
	BuiltinEventHandler
}

func (mh MyHandler) OnConn(c Conn) error{
	p, err := c.Next(-1)
    if err!=nil{
        return evnet.ErrClose
    }
    c.AsyncWrite(p, func(c Conn){
        c.Close()
    })
}

//main.go
func main() {
	readHandler := MyHandler{BuiltinEventHandler{}}
	Run(readHandler, "tcp://127.0.0.1:9000")
}
```

## BenchTest

```
# Machine information
        OS : Ubuntu 20.04/x86_64
       CPU : 8 CPU cores, Intel(R) Core(TM) i5-8265U CPU @ 1.60GHz
    Memory : 4.0 GiB
  Platform : VMware Workstation 16.2.3

# Go version and settings
Go Version : go1.18.3 linux/amd64
GOMAXPROCS : 8
```

```
The average performance of evnet in the condition of 12 threads and 1000 connections is 103% more than standard net/http package.

| Thread Stats | Avg    | Stdev  | Max     | +/- Stdev |
| ------------ | ------ | ------ | ------- | --------- |
| Latency      | 7.94ms | 8.90ms | 84.67ms | 85.92%    |
| Req/Sec      | 13.60k | 4.41k  | 32.25k  | 68.19%    |
1622016 requests in 10.09s, 199.55MB read

Requests/sec: 160690.88
Transfer/sec:     19.77MB
```