package main

import (
	"esl-service/internal/app"
	"flag"
	"runtime"
)

var (
	fshost   = flag.String("fshost", "host.docker.internal", "Freeswitch hostname. Default: localhost")
	fsport   = flag.Uint("fsport", 8021, "Freeswitch port. Default: 8021")
	password = flag.String("pass", "ClueCon", "Freeswitch password. Default: ClueCon")
	timeout  = flag.Int("timeout", 10, "Freeswitch conneciton timeout in seconds. Default: 10")
)

func main() {
	// Boost it as much as it can go ...
	runtime.GOMAXPROCS(runtime.NumCPU())

	client, err := app.CreateClient(*fshost, *fsport, *password, *timeout)
	if err != nil {
		return
	}
	defer client.Close()

	for {
		msg, err := app.GetMessage(client)

		if err != nil {
			break
		}

		app.Execute(client, msg)
	}
}
