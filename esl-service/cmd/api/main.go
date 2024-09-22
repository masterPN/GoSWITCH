package main

import (
	"esl-service/internal/app"
	"flag"
	"runtime"
	"strings"

	"github.com/0x19/goesl"
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

	client, err := goesl.NewClient(*fshost, *fsport, *password, *timeout)

	if err != nil {
		goesl.Error("Error while creating new client: %s", err)
		return
	}

	goesl.Debug("Yuhu! New client: %q", client)

	// Apparently all is good... Let us now handle connection :)
	// We don't want this to be inside of new connection as who knows where it my lead us.
	// Remember that this is crutial part in handling incoming messages :)
	go client.Handle()

	client.Send("events json ALL")

	for {
		msg, err := client.ReadMessage()

		if err != nil {

			// If it contains EOF, we really dont care...
			if !strings.Contains(err.Error(), "EOF") && err.Error() != "unexpected end of JSON input" {
				goesl.Error("Error while reading Freeswitch message: %s", err)
			}
			break
		}

		app.Execute(client, msg.Headers)
	}

}
