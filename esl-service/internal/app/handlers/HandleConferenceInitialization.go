package handlers

import (
	"esl-service/internal/app/helpers"
	"log"
	"strings"

	"github.com/0x19/goesl"
)

func HandleConferenceInitialization(client *goesl.Client, msg map[string]string) {
	sipPort, externalDomain := helpers.LoadConfiguration()

	initConferenceData := strings.Split(msg["variable_current_application_data"], ", ")

	err := helpers.InitiateConferenceCalls(client, initConferenceData, externalDomain, sipPort, msg)
	if err != nil {
		log.Printf("Error originating calls: %s\n", err)
	}
}
