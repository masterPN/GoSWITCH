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
		client.BgApi(fmt.Sprintf("originate {origination_caller_id_number=%s}sofia/internal/%s:%v &conference(%s)",
			initConferenceData[2], initConferenceData[3], sipPort,
			initConferenceData[2]))
	}
}