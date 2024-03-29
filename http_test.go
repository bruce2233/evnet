package evnet

import (
	"net"
	"net/http"
	"os"
	"runtime/pprof"
	"strconv"
	"testing"
	"time"

	"github.com/evanphx/wildcat"
	log "github.com/sirupsen/logrus"
)

type HttpServer struct {
	BuiltinEventHandler
}

type httpCodec struct {
	parser *wildcat.HTTPParser
	buf    []byte
}

func (hs *HttpServer) OnBoot(mr *MainReactor) error {
	log.Infof("\n=================Welcome!=================")
	log.Infof("\n███████╗██╗░░░██╗███╗░░██╗███████╗████████╗\n██╔════╝██║░░░██║████╗░██║██╔════╝╚══██╔══╝\n█████╗░░╚██╗░██╔╝██╔██╗██║█████╗░░░░░██║░░░\n██╔══╝░░░╚████╔╝░██║╚████║██╔══╝░░░░░██║░░░\n███████╗░░╚██╔╝░░██║░╚███║███████╗░░░██║░░░\n╚══════╝░░░╚═╝░░░╚═╝░░╚══╝╚══════╝░░░╚═╝░░░")
	return nil
}

func (hs *HttpServer) OnTraffic(c Conn) error {
	// httpParser := wildcat.HTTPParser{}
	// println(httpParser.Version)
	hc := c.Context().(*httpCodec)
	buf, _ := c.Next(-1)
	// log.Println(hc, buf)
	_, err := hc.parser.Parse(buf)
	if err != nil {
		c.Write([]byte("Internal Server Error"))
		return ErrClose
	}
	if string(hc.parser.Path) == "/close" {
		return ErrShutdown
	}
	hc.appendResponse()
	bodyLen := int(hc.parser.ContentLength())
	if bodyLen == -1 {
		bodyLen = 0
	}
	// buf = buf[headerOffset+bodyLen:]
	hc.buf = append(hc.buf, "Content-Length: "+strconv.Itoa(len(hc.parser.Path))+"\r\n\r\n"+string(hc.parser.Path)...)
	println(string(hc.buf))
	// hc.buf = append(hc.buf, "Content-Length: 12\r\n\r\nHello World!"...)
	c.AsyncWrite(hc.buf, func(c Conn) error {
		// c.Close()

		return nil
	})
	hc.buf = hc.buf[:0]
	return nil
}

func (hc *httpCodec) appendResponse() {
	hc.buf = append(hc.buf, "HTTP/1.1 200 OK\r\nServer: evnet\r\nContent-Type: text/plain\r\nDate: "...)
	hc.buf = time.Now().AppendFormat(hc.buf, "Mon, 02 Jan 2006 15:04:05 GMT\r\n")
	// hc.buf = append(hc.buf, "Content-Length: 12\r\n\r\nHello World!"...)
}

func (hs *HttpServer) OnOpen(c Conn) error {
	httpParser := wildcat.NewHTTPParser()

	c.SetContext(&httpCodec{
		parser: httpParser,
	})
	log.Println("OnOpen: httpCodec ", c.Context().(*httpCodec).parser.TotalHeaders)
	return nil
}

func TestHttpBench(t *testing.T) {
	f, _ := os.OpenFile("cpu.pprof", os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	defer println("defer works")
	hs := new(HttpServer)
	Run(hs, "tcp://127.0.0.1:9000", WithLogLevel(log.InfoLevel))
}

func TestWriteHTTPRequest(t *testing.T) {
	resp, err := http.Get("http://127.0.0.1:9000/index.html")
	if err != nil {
		t.Log(err)
	}
	t.Log(resp)
}

func TestTCPWrite(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	defer conn.Close()
	if err != nil {
		t.Log(err)
	}
	b := make([]byte, 1024)
	conn.Write([]byte("Hello World!"))
	n, err := conn.Read(b)
	if err != nil {
		t.Log(err)
	}
	println(b[:n])
}

func TestWildCatParse(t *testing.T) {
	httpParser := wildcat.NewHTTPParser()
	// str := "GET /example HTTP/1.1\r\nHost: 192.168.87.141:9000"
	str := "GET /example HTTP/1.1\r\nHost: 192.168.87.141:9000\r\n\r\n"
	// str := "GET /example HTTP/1.1\r\nHost: 192.168.87.141:90"
	buf := []byte(str)
	t.Log(str, len(buf))
	headerOffset, err := httpParser.Parse(buf)
	t.Log(headerOffset, err, httpParser.ContentLength())
}
