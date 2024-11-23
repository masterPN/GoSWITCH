package handlers

import (
	"bytes"
	"encoding/json"
	"esl-service/internal/app/helpers"
	"esl-service/internal/data"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/0x19/goesl"
)

// Configuration loading function
func loadConfiguration() (int, string) {
	sipPort, _ := strconv.Atoi(os.Getenv("SIP_PORT"))
	externalDomain := os.Getenv("EXTERNAL_DOMAIN")
	return sipPort, externalDomain
}

// Main Conference Call handler
func initiateConferenceCalls(client *goesl.Client, initConferenceData []string, externalDomain string, sipPort int, msg map[string]string) error {
	if validateRadiusAndHandleConference(client, initConferenceData, msg) {
		return nil
	}

	routingResponse, err := fetchOperatorRouting(helpers.NormalizeDestinationNumber(initConferenceData[3]))
	if err != nil {
		log.Printf("Error fetching operator routing: %s\n", err)
		return err
	}

	baseClassResponse := []int{
		routingResponse.Class1,
		routingResponse.Class2,
		routingResponse.Class3,
	}

	for _, response := range baseClassResponse {
		if response == 0 {
			continue
		}

		operatorMapping, err := fetchInternalCodemapping(strconv.Itoa(response))
		if err != nil {
			log.Printf("Error fetching internal code mapping: %s\n", err)
			continue
		}

		operatorPrefix := operatorMapping.OperatorCode
		if originateCall(client, initConferenceData, response, strconv.Itoa(operatorPrefix), externalDomain, sipPort) {
			return nil
		}
	}

	goesl.Debug("There's no operator available.")
	client.BgApi(fmt.Sprintf(destroyConferenceCommand, initConferenceData[1]))
	http.Get(fmt.Sprintf("http://redis-service:8080/popRadiusAccountingData/%s", initConferenceData[1]))
	return nil
}

// Helper functions for network communication

// Validate radius and handle conference destruction if necessary
func validateRadiusAndHandleConference(client *goesl.Client, conferenceInitData []string, msg map[string]string) bool {
	radiusValidationRequest, _ := json.Marshal(map[string]string{
		"prefix":            strings.Replace(msg[callerDestinationHeader], conferenceInitData[3], "", 1),
		"callingNumber":     helpers.NormalizeDestinationNumber(conferenceInitData[2]),
		"destinationNumber": helpers.NormalizeDestinationNumber(conferenceInitData[3]),
	})

	radiusValidationResponse, err := http.Post("http://mssql-service:8080/radiusOnestageValidate", jsonContentType, bytes.NewBuffer(radiusValidationRequest))
	if err != nil {
		log.Printf("Error validating radius: %s\n", err)
		return true
	}
	defer radiusValidationResponse.Body.Close()

	var radiusValidationResponseData data.RadiusOnestageValidateData
	if err := json.NewDecoder(radiusValidationResponse.Body).Decode(&radiusValidationResponseData); err != nil {
		log.Printf("Error decoding radius response: %s\n", err)
		return true
	}

	if radiusValidationResponseData.Status > 2 {
		client.BgApi(fmt.Sprintf(destroyConferenceCommand, conferenceInitData[1]))
		log.Printf("Kicked due to radius status %v\n", radiusValidationResponseData.Status)
		return true
	}

	radiusAccountingData := data.RadiusAccounting{
		AccessNo:     radiusValidationResponseData.PrefixNo,
		Anino:        helpers.NormalizeDestinationNumber(conferenceInitData[2]),
		DestNo:       helpers.NormalizeDestinationNumber(conferenceInitData[3]),
		SubscriberNo: radiusValidationResponseData.AccountNum,
		SessionID:    conferenceInitData[1],
		InTrunkID:    0,
		ReasonID:     radiusValidationResponseData.Status,
	}

	radiusAccountingRequest, _ := json.Marshal(radiusAccountingData)

	_, err = http.Post("http://redis-service:8080/saveRadiusAccountingData", jsonContentType, bytes.NewBuffer(radiusAccountingRequest))
	if err != nil {
		log.Printf("POST http://redis-service:8080/saveRadiusAccountingData - %s\n", err)
		return true
	}

	return false
}

// Fetch operator routing for a given destination
func fetchOperatorRouting(destination string) (data.OptimalRouteData, error) {
	resp, err := http.Get(fmt.Sprintf("http://mssql-service:8080/optimalRoute?pCallString=%s", destination))
	if err != nil {
		return data.OptimalRouteData{}, err
	}
	defer resp.Body.Close()

	var routingResponse data.OptimalRouteData
	if err := json.NewDecoder(resp.Body).Decode(&routingResponse); err != nil {
		return data.OptimalRouteData{}, err
	}
	return routingResponse, nil
}

// Fetch the route data for a given route
func fetchInternalCodemapping(internalCode string) (data.InternalCodemappingData, error) {
	resp, err := http.Get(fmt.Sprintf("http://mssql-service:8080/internalCodemapping/internalCode=%s", internalCode))
	if err != nil {
		return data.InternalCodemappingData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		var internalCodemappingResponse struct {
			Error                   string                       `json:"error"`
			InternalCodemappingData data.InternalCodemappingData `json:"internalCodemappingData"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&internalCodemappingResponse); err != nil {
			return data.InternalCodemappingData{}, err
		}

		return internalCodemappingResponse.InternalCodemappingData, fmt.Errorf(internalCodemappingResponse.Error)
	}

	var internalCodemappingResponse data.InternalCodemappingData
	if err := json.NewDecoder(resp.Body).Decode(&internalCodemappingResponse); resp.StatusCode == http.StatusOK && err != nil {
		return data.InternalCodemappingData{}, err
	}
	return internalCodemappingResponse, nil
}

// Originate call to a given operator

func originateCall(client *goesl.Client, initConferenceData []string, baseClass int, operatorPrefix, externalDomain string, sipPort int) bool {
	client.BgApi(fmt.Sprintf("originate {origination_caller_id_number=%s}sofia/external/%s%s@%s:%v &conference(%s)",
		initConferenceData[2], operatorPrefix, helpers.NormalizeDestinationNumber(initConferenceData[3]), externalDomain, sipPort,
		initConferenceData[1]))

	return waitForCall(client, baseClass, operatorPrefix, helpers.NormalizeDestinationNumber(initConferenceData[3]), initConferenceData[1])
}

// Wait for the call to be established and handle various states
func waitForCall(client *goesl.Client, baseClass int, operatorPrefix, destination, conferenceName string) bool {
	secondaryClient, err := helpers.CreateClient()
	if err != nil {
		goesl.Debug("Create secondary client in waitForCall failed!")
		return false
	}
	defer secondaryClient.Close()

	startTime := time.Now()
	var once sync.Once

	for {
		if !withinTimeout(startTime) {
			break
		}

		handleBgApiCall(client, startTime, &once)

		msg, err := secondaryClient.ReadMessage()
		if err != nil {
			handleReadError(err)
			break
		}

		if handleCalleeAndConnected(msg, client, baseClass, operatorPrefix, destination, conferenceName) {
			return true
		}

		if isOperatorUnavailable(msg, operatorPrefix, destination) {
			logOperatorIssue(msg, operatorPrefix)
			return false
		}
	}

	goesl.Debug("WARNING - There's no matched case for %q", operatorPrefix+destination)
	return false
}

// Helper functions for processing call results

func withinTimeout(startTime time.Time) bool {
	return time.Since(startTime) <= 5*time.Second
}

func handleBgApiCall(client *goesl.Client, startTime time.Time, once *sync.Once) {
	if time.Since(startTime) > 2*time.Second {
		once.Do(func() {
			client.BgApi("show channels")
		})
	}
}

func handleCalleeAndConnected(msg *goesl.Message, client *goesl.Client, baseClass int, operatorPrefix, destination, conferenceName string) bool {
	if isConnected(msg, operatorPrefix, destination) {
		return handleConnectedCall(client, baseClass, conferenceName)
	}

	if isCalleeUnavailable(msg, operatorPrefix, destination) {
		handleCalleeIssue(client, msg, operatorPrefix, destination, conferenceName)
		return true
	}

	return false
}

// Helper functions for managing specific call states

func handleConnectedCall(client *goesl.Client, baseClass int, conferenceName string) bool {
	goesl.Debug("Received call, exiting initConferenceHandler")

	radiusAccountingBody := data.RadiusAccounting{
		SessionID:  conferenceName,
		OutTrunkID: baseClass,
	}

	postBody, _ := json.Marshal(radiusAccountingBody)

	_, err := http.Post("http://redis-service:8080/saveRadiusAccountingData", jsonContentType, bytes.NewBuffer(postBody))
	if err != nil {
		log.Printf("POST http://redis-service:8080/saveRadiusAccountingData - %s\n", err)
		client.BgApi(fmt.Sprintf(destroyConferenceCommand, conferenceName))
		return true
	}

	return true
}

func handleCalleeIssue(client *goesl.Client, msg *goesl.Message, operatorPrefix, destination, conferenceName string) {
	logCalleeIssue(msg, operatorPrefix, destination)
	client.BgApi(fmt.Sprintf(destroyConferenceCommand, conferenceName))
	http.Get(fmt.Sprintf("http://redis-service:8080/popRadiusAccountingData/%s", conferenceName))
}

func handleReadError(err error) {
	if !strings.Contains(err.Error(), "EOF") && err.Error() != "unexpected end of JSON input" {
		goesl.Error("Error while reading Freeswitch message: %s", err)
	}
}

func isConnected(msg *goesl.Message, operatorPrefix, destination string) bool {
	return (msg.Headers["Action"] == "add-member" &&
		msg.Headers[answerStateHeader] == "early" &&
		msg.Headers[callerDestinationHeader] == operatorPrefix+destination) ||
		(msg.Headers["Event-Calling-Function"] == "bgapi_exec" &&
			msg.Headers["Job-Command"] == "show" &&
			msg.Headers["Job-Command-Arg"] == "channels" &&
			strings.Contains(msg.Headers["body"], operatorPrefix+destination+",conference"))
}

func isCalleeUnavailable(msg *goesl.Message, operatorPrefix, destination string) bool {
	return msg.Headers[answerStateHeader] == "hangup" &&
		msg.Headers[callerDestinationHeader] == operatorPrefix+destination &&
		slices.Contains(sipCalleeUnavailableCode, msg.Headers["variable_hangup_cause_q850"])
}

func logCalleeIssue(msg *goesl.Message, operatorPrefix, destination string) {
	goesl.Debug(`%q has a problem, please contact callee %q.\n
		code - %q, reason - %q`,
		operatorPrefix+destination, destination,
		msg.Headers["variable_hangup_cause_q850"], msg.Headers["variable_sip_invite_failure_phrase"])
}

func isOperatorUnavailable(msg *goesl.Message, operatorPrefix, destination string) bool {
	return msg.Headers[answerStateHeader] == "hangup" &&
		msg.Headers[callerDestinationHeader] == operatorPrefix+destination &&
		slices.Contains(sipOperatorUnavailableCode, msg.Headers["variable_hangup_cause_q850"])
}

func logOperatorIssue(msg *goesl.Message, operatorPrefix string) {
	goesl.Debug(`%q has a problem, please contact operator %q.\n
		code - %q, reason - %q`,
		operatorPrefix, operatorPrefix,
		msg.Headers["variable_hangup_cause_q850"], msg.Headers["variable_sip_invite_failure_phrase"])
}
