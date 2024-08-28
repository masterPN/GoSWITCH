package app

import (
	"esl-service/internal/app/handlers"
	"strings"

	. "github.com/0x19/goesl"
)

const (
	eventCallingFunction = "Event-Calling-Function"
)

func Execute(client *Client, msg map[string]string) {
	switch {
	case strings.Contains(msg["variable_current_application_data"], "initConference"):
		go handlers.InitConferenceHandler(client, msg)
	case strings.Contains(msg["Hangup-Cause"], "CALL_REJECTED") &&
		strings.Contains(msg[eventCallingFunction], "switch_channel_perform_hangup"):
		// Callee rejects the call
		go handlers.RejectConferenceHandler(client, msg)
	case strings.Contains(msg["Answer-State"], "answered") &&
		strings.Contains(msg["Call-Direction"], "outbound") &&
		strings.Contains(msg[eventCallingFunction], "switch_channel_perform_mark_answered"):
		// Callee accepts the call
		go handlers.JoinConferenceHandler(msg)
	case strings.Contains(msg["Answer-State"], "hangup") &&
		strings.Contains(msg["Call-Direction"], "inbound") &&
		strings.Contains(msg[eventCallingFunction], "switch_core_session_perform_destroy"):
		go handlers.EndConferenceHandler(client, msg)
	}
}
