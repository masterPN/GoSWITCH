package app

import (
	"esl-service/internal/app/handlers"
	"strings"

	"github.com/0x19/goesl"
)

const (
	VariableCurrentApplicationData = "variable_current_application_data"
	EventName                      = "Event-Name"
	HangupCause                    = "Hangup-Cause"
	AnswerState                    = "Answer-State"
	CallDirection                  = "Call-Direction"
	EventCallingFunction           = "Event-Calling-Function"
)

func Execute(client *goesl.Client, msg map[string]string) {
	eventName := msg[EventName]
	eventFunction := msg[EventCallingFunction]
	applicationData := msg[VariableCurrentApplicationData]
	hangupCause := msg[HangupCause]
	answerState := msg[AnswerState]
	callDirection := msg[CallDirection]

	switch {
	case strings.Contains(applicationData, "initConference") && eventName == "CHANNEL_EXECUTE_COMPLETE":
		goesl.Debug("initConference\n %q", applicationData)
		go handlers.InitConferenceHandler(client, msg)
	case strings.Contains(hangupCause, "CALL_REJECTED") && strings.Contains(eventFunction, "switch_channel_perform_hangup"):
		// Callee rejects the call
		goesl.Debug("callee rejected\n %q", msg)
		go handlers.RejectConferenceHandler(client, msg)
	case strings.Contains(answerState, "answered") &&
		strings.Contains(callDirection, "outbound") &&
		strings.Contains(eventFunction, "switch_channel_perform_mark_answered"):
		// Callee accepts the call
		goesl.Debug("callee accepted\n %q", msg)
		go handlers.JoinConferenceHandler(msg)
	case strings.Contains(answerState, "hangup") &&
		strings.Contains(callDirection, "inbound") &&
		strings.Contains(eventFunction, "switch_core_session_perform_destroy"):
		goesl.Debug("end conference\n %q", msg)
		go handlers.EndConferenceHandler(client, msg)
	}
}
