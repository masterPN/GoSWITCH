package helpers

import (
	"flag"
	"strings"

	"github.com/0x19/goesl"
)

var (
	fshost   = flag.String("fshost", "host.docker.internal", "Freeswitch hostname. Default: localhost")
	fsport   = flag.Uint("fsport", 8021, "Freeswitch port. Default: 8021")
	password = flag.String("pass", "ClueCon", "Freeswitch password. Default: ClueCon")
	timeout  = flag.Int("timeout", 10, "Freeswitch conneciton timeout in seconds. Default: 10")
)

func CreateClient() (*goesl.Client, error) {
	client, err := goesl.NewClient(*fshost, *fsport, *password, *timeout)
	if err != nil {
		goesl.Error("Error while creating new client: %s", err)
		return nil, err
	}

	goesl.Debug("Yuhu! New client: %q", client)

	go client.Handle()
	client.Send("events json ALL")

	return client, nil
}

func GetMessage(client *goesl.Client) (map[string]string, error) {
	msg, err := client.ReadMessage()
	if err != nil {
		if shouldLogError(err) {
			goesl.Error("Error while reading Freeswitch message: %s", err)
		}
		return nil, err
	}
	return msg.Headers, nil
}

func shouldLogError(err error) bool {
	errMsg := err.Error()
	return !strings.Contains(errMsg, "EOF") && errMsg != "unexpected end of JSON input"
}
