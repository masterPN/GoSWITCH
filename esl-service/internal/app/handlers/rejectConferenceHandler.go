package handlers

import (
	"fmt"
	"net/http"

	"github.com/0x19/goesl"
)

func RejectConferenceHandler(client *goesl.Client, msg map[string]string) {
	client.BgApi(fmt.Sprintf("conference %v kick all", msg["variable_conference_name"]))
	http.Get(fmt.Sprintf("http://redis-service:8080/popRadiusAccountingData/%s", msg["variable_conference_name"]))
}
