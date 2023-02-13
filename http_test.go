package evnet

import (
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/evanphx/wildcat"
)

type HttpServer struct {
	BuiltinEventHandler
}

type httpCodec struct {
	parser *wildcat.HTTPParser
	buf    []byte
}

func (hs *HttpServer) OnBoot(mr *MainReactor) error {
	log.Println("\n=================Welcome!=================")
	log.Println("\n███████╗██╗░░░██╗███╗░░██╗███████╗████████╗\n██╔════╝██║░░░██║████╗░██║██╔════╝╚══██╔══╝\n█████╗░░╚██╗░██╔╝██╔██╗██║█████╗░░░░░██║░░░\n██╔══╝░░░╚████╔╝░██║╚████║██╔══╝░░░░░██║░░░\n███████╗░░╚██╔╝░░██║░╚███║███████╗░░░██║░░░\n╚══════╝░░░╚═╝░░░╚═╝░░╚══╝╚══════╝░░░╚═╝░░░")
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
	hc.appendResponse()
	bodyLen := int(hc.parser.ContentLength())
	if bodyLen == -1 {
		bodyLen = 0
	}
	// buf = buf[headerOffset+bodyLen:]

	c.AsyncWrite(hc.buf, nil)
	hc.buf = hc.buf[:0]
	return nil
}

func (hc *httpCodec) appendResponse() {
	hc.buf = append(hc.buf, "HTTP/1.1 200 OK\r\nServer: evnet\r\nContent-Type: text/plain\r\nDate: "...)
	hc.buf = time.Now().AppendFormat(hc.buf, "Mon, 02 Jan 2006 15:04:05 GMT")
	hc.buf = append(hc.buf, "\r\nContent-Length: 12\r\n\r\nHello World!"...)
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
	hs := new(HttpServer)
	Run(hs, "tcp://192.168.87.141:9000")
}

func TestWriteHTTPRequest(t *testing.T) {
	resp, err := http.Get("http://192.168.87.141:9000/example")
	if err != nil {
		t.Log(err)
	}
	t.Log(resp)
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
