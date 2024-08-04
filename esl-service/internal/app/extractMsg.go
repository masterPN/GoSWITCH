package app

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	. "github.com/0x19/goesl"
)

func Execute(client *Client, msg map[string]string) {
	sipPort, _ := strconv.Atoi(os.Getenv("SIP_PORT"))

	if strings.Contains(msg["variable_current_application_data"], "initConference") {
		// Call is starting.
		slog.Info(msg["variable_current_application_data"])

		initConferenceData := strings.Split(msg["variable_current_application_data"], ", ")
		fmt.Printf("originate sofia/internal/%s:%v &conference(%s)",
			initConferenceData[3], sipPort,
			initConferenceData[2])
		client.BgApi(fmt.Sprintf("originate sofia/internal/%s:%v &conference(%s)",
			initConferenceData[3], sipPort,
			initConferenceData[2]))
	}
}
