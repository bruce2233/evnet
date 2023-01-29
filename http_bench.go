package evnet
import(
	"testing"

	"github.com/evanphx/wildcat"
)

type HttpServer struct{
	BuiltinEventHandler
}

func (hs HttpServer) OnTraffic(c Conn){
	httpParser := wildcat.HTTPParser{}
	println(httpParser.Version)
}

func TestHttpBench(t *testing.T){

}

