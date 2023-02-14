package evnet

import (
	log "github.com/sirupsen/logrus"

	"os"
	"strings"
)

func Run(eventHandler EventHandler, protoAddr string, optList ...SetOption) {
	opts := loadOptions(optList)
	if opts.LogPath != "" {
		outputFile, err := os.OpenFile(opts.LogPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
		if err != nil {
			log.Error("log file open error")
		}
		log.SetOutput(outputFile)
	}
	// log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	//options working
	network, address := parseProtoAddr(protoAddr)

	//mainReactor initialize
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
	(**mr.eventHandlerPP).OnBoot(mr)
	go mr.Loop()
	mr.waitForShutdown()
	mr.stop()
}
