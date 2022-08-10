package reactor

import "strings"

func Run(eventHandler EventHandler, protoAddr string) {
	//options working
	network, address := parseProtoAddr(protoAddr)

	mr := new(MainReactor)
	mr.Init(network, address)
	mr.SetEventHandler(eventHandler)
	serve(mr)
}

func parseProtoAddr(protoAddr string) (network, address string) {
	network = "tcp"
	address = strings.ToLower(protoAddr)
	if strings.Contains(address, "://") {
		pair := strings.Split(address, "://")
		network = pair[0]
		address = pair[1]
	}
	return
}

func serve(mr *MainReactor) {
	mr.Loop()
}
