package evnet

import (
	"log"
	"net/http"
	"testing"

	"github.com/evanphx/wildcat"
)

type HttpServer struct {
	BuiltinEventHandler
}

type httpCodec struct {
	parser *wildcat.HTTPParser
	buf    []byte
}

func (hs *HttpServer) OnTraffic(c Conn) error {
	// httpParser := wildcat.HTTPParser{}
	// println(httpParser.Version)
	return ErrShutdown
}

func (hs *HttpServer) OnOpen(c Conn) error {
	httpParser := wildcat.NewHTTPParser()

	c.SetContext(&httpCodec{
		parser: httpParser,
	})
	log.Println("OnOpen: httpCodec ", c.Context().(*httpCodec).parser.TotalHeaders)
	return ErrClose
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
