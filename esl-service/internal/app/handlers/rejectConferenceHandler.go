package handlers

import (
	"fmt"

	"github.com/0x19/goesl"
)

func RejectConferenceHandler(client *goesl.Client, msg map[string]string) {
	client.BgApi(fmt.Sprintf("conference %v kick all", msg["Caller-Caller-ID-Number"]))
}
