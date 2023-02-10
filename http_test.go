package evnet

import (
	"log"
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

func (hs *HttpServer) OnTraffic(c Conn) {
	// httpParser := wildcat.HTTPParser{}
	// println(httpParser.Version)
}

func (hs *HttpServer) OnOpen(c Conn) {
	httpParser := wildcat.NewHTTPParser()

	c.SetContext(&httpCodec{
		parser: httpParser,
	})
	log.Println("OnOpen: httpCodec ", c.Context().(*httpCodec).parser.TotalHeaders)
}

func TestHttpBench(t *testing.T) {
	hs := new(HttpServer)
	Run(hs, "tcp://192.168.87.141:9000")
}
