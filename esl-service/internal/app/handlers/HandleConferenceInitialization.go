package handlers

import (
	"log"
	"strings"

	"github.com/0x19/goesl"
)

const (
	answerStateHeader        = "Answer-State"
	callerDestinationHeader  = "Caller-Destination-Number"
	destroyConferenceCommand = "conference %v kick all"
	jsonContentType          = "application/json"
)

var (
	sipOperatorUnavailableCode = []string{"1", "41"}
	sipCalleeUnavailableCode   = []string{"17"}
)

func HandleConferenceInitialization(client *goesl.Client, msg map[string]string) {
	sipPort, externalDomain := loadConfiguration()

	initConferenceData := strings.Split(msg["variable_current_application_data"], ", ")

	err := initiateConferenceCalls(client, initConferenceData, externalDomain, sipPort, msg)
	if err != nil {
		log.Printf("Error originating calls: %s\n", err)
	}
}
