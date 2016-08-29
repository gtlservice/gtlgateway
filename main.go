package main

import "github.com/gtlservice/gtlgateway/base"
import "github.com/gtlservice/gtlgateway/server"
import "github.com/gtlservice/gutils/system"

import (
	"os"
)

func main() {

	gateway, err := server.NewGateway()
	if err != nil {
		panic(err)
		os.Exit(base.EXITCODE_INITFAILED)
	}

	defer func() {
		gateway.Stop()
		os.Exit(base.EXITCODE_EXITED)
	}()

	if err := gateway.Startup(); err != nil {
		panic(err)
		os.Exit(base.EXITCODE_STARTFAILED)
	}
	system.InitSignal(nil)
}
