package helpers

import (
	"bytes"
	"encoding/json"
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
	sharedData "github.com/masterPN/GoSWITCH-shared/data"
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

// Configuration loading function
func LoadConfiguration() (int, string) {
	sipPort, _ := strconv.Atoi(os.Getenv("SIP_PORT"))
	externalDomain := os.Getenv("EXTERNAL_DOMAIN")
	return sipPort, externalDomain
}

// Main Conference Call handler
func InitiateConferenceCalls(client *goesl.Client, conferenceData []string, externalDomain string, sipPort int, message map[string]string) error {
	if ValidateRadiusAndHandleConference(client, conferenceData, message) {
		return nil
	}

	routeData, err := FetchOptimalRouteData(NormalizeDestinationNumber(conferenceData[3]))
	if err != nil {
		log.Printf("Error fetching operator routing: %s\n", err)
		return err
	}

	baseClasses := []int{routeData.Class1, routeData.Class2, routeData.Class3}

	for _, baseClass := range baseClasses {
		if baseClass == 0 {
			continue
		}

		mappingData, err := FetchInternalCodemapping(strconv.Itoa(baseClass))
		if err != nil {
			log.Printf("Error fetching internal code mapping: %s\n", err)
			continue
		}

		operatorPrefix := strconv.Itoa(mappingData.OperatorCode)
		if OriginateCallToOperator(client, conferenceData, baseClass, operatorPrefix, externalDomain, sipPort) {
			return nil
		}
	}

	log.Println("No operator available.")
	client.BgApi(fmt.Sprintf(destroyConferenceCommand, conferenceData[1]))
	http.Get(fmt.Sprintf("http://redis-service:8080/popRadiusAccountingData/%s", conferenceData[1]))
	return nil
}

// Helper functions for network communication

// Validate radius and handle conference destruction if necessary
func ValidateRadiusAndHandleConference(client *goesl.Client, conferenceInitData []string, msg map[string]string) bool {
	requestData := map[string]string{
		"prefix":            strings.Replace(msg[callerDestinationHeader], conferenceInitData[3], "", 1),
		"callingNumber":     NormalizeDestinationNumber(conferenceInitData[2]),
		"destinationNumber": NormalizeDestinationNumber(conferenceInitData[3]),
	}

	requestBytes, _ := json.Marshal(requestData)
	response, err := http.Post("http://mssql-service:8080/radiusOnestageValidate", jsonContentType, bytes.NewBuffer(requestBytes))
	if err != nil {
		return true
	}
	defer response.Body.Close()

	var responseData sharedData.RadiusOnestageValidateData
	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return true
	}

	if responseData.Status > 2 {
		client.BgApi(fmt.Sprintf(destroyConferenceCommand, conferenceInitData[1]))
		return true
	}

	accountingData := data.RadiusAccounting{
		AccessNo:     responseData.PrefixNo,
		Anino:        NormalizeDestinationNumber(conferenceInitData[2]),
		DestNo:       NormalizeDestinationNumber(conferenceInitData[3]),
		SubscriberNo: responseData.AccountNum,
		SessionID:    conferenceInitData[1],
		InTrunkID:    0,
		ReasonID:     responseData.Status,
	}

	accountingRequest, _ := json.Marshal(accountingData)

	_, err = http.Post("http://redis-service:8080/saveRadiusAccountingData", jsonContentType, bytes.NewBuffer(accountingRequest))

	return err != nil
}

// FetchOptimalRouteData fetches the optimal route data for a given destination
func FetchOptimalRouteData(destination string) (data.OptimalRouteData, error) {
	url := fmt.Sprintf("http://mssql-service:8080/optimalRoute?pCallString=%s", destination)
	resp, err := http.Get(url)
	if err != nil {
		return data.OptimalRouteData{}, err
	}
	defer resp.Body.Close()

	var optimalRouteData data.OptimalRouteData
	if err := json.NewDecoder(resp.Body).Decode(&optimalRouteData); err != nil {
		return data.OptimalRouteData{}, err
	}
	return optimalRouteData, nil
}

// FetchInternalCodemapping fetches the internal code mapping for a given internal code
func FetchInternalCodemapping(internalCode string) (data.InternalCodemappingData, error) {
	resp, err := http.Get(fmt.Sprintf("http://redis-service:8080/internalCodemappingData/%s", internalCode))
	if err != nil {
		return data.InternalCodemappingData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		var internalCodemappingError data.InternalCodemappingDataError

		if err := json.NewDecoder(resp.Body).Decode(&internalCodemappingError); err != nil {
			return data.InternalCodemappingData{}, err
		}

		return data.InternalCodemappingData{}, fmt.Errorf(internalCodemappingError.Error)
	}

	var internalCodemappingData data.InternalCodemappingData
	if err := json.NewDecoder(resp.Body).Decode(&internalCodemappingData); resp.StatusCode == http.StatusOK && err != nil {
		return data.InternalCodemappingData{}, err
	}

	return internalCodemappingData, nil
}

// Originate call to a given operator

func OriginateCallToOperator(client *goesl.Client, conferenceData []string, operatorCode int, operatorPrefix, externalDomain string, sipPort int) bool {
	callerIDNumber := conferenceData[2]
	calleeNumber := NormalizeDestinationNumber(conferenceData[3])
	conferenceName := conferenceData[1]

	client.BgApi(fmt.Sprintf(
		"originate {origination_caller_id_number=%s}sofia/external/%s%s@%s:%v &conference(%s)",
		callerIDNumber, operatorPrefix, calleeNumber, externalDomain, sipPort, conferenceName,
	))

	return waitForCallToBeEstablished(client, operatorCode, operatorPrefix, calleeNumber, conferenceName)
}

// Wait for the call to be established and handle various states
func waitForCallToBeEstablished(client *goesl.Client, baseClass int, operatorPrefix, destination, conferenceName string) bool {
	secondaryClient, err := CreateClient()
	if err != nil {
		log.Printf("Failed to create secondary client in waitForCall: %v", err)
		return false
	}
	defer secondaryClient.Close()

	startTime := time.Now()
	var once sync.Once

	for {
		if !isWithinTimeout(startTime) {
			break
		}

		handleBackgroundApiCall(client, startTime, &once)

		message, err := secondaryClient.ReadMessage()
		if err != nil {
			handleReadError(err)
			break
		}

		if processCalleeAndConnection(message, client, baseClass, operatorPrefix, destination, conferenceName) {
			return true
		}

		if isOperatorUnavailable(message, operatorPrefix, destination) {
			logOperatorIssue(message, operatorPrefix)
			return false
		}
	}

	log.Printf("No matching case for operator and destination: %s", operatorPrefix+destination)
	return false
}

// Helper functions for processing call results

func isWithinTimeout(start time.Time) bool {
	return time.Since(start) <= 5*time.Second
}

func handleBackgroundApiCall(client *goesl.Client, startTime time.Time, once *sync.Once) {
	if time.Since(startTime) > 2*time.Second {
		once.Do(func() {
			client.BgApi("show channels")
		})
	}
}

func processCalleeAndConnection(message *goesl.Message, client *goesl.Client, baseClass int, operatorPrefix, destination, conferenceName string) bool {
	if isConnected(message, operatorPrefix, destination) {
		return handleConnectedCall(client, baseClass, conferenceName)
	}

	if isCalleeUnavailable(message, operatorPrefix, destination) {
		handleCalleeIssue(client, message, operatorPrefix, destination, conferenceName)
		return true
	}

	return false
}

// Helper functions for managing specific call states

func handleConnectedCall(client *goesl.Client, baseClass int, conferenceName string) bool {
	radiusAccountingData := data.RadiusAccounting{
		SessionID:  conferenceName,
		OutTrunkID: baseClass,
	}

	postBody, err := json.Marshal(radiusAccountingData)
	if err != nil {
		log.Printf("Error marshaling radius accounting data: %s\n", err)
		client.BgApi(fmt.Sprintf(destroyConferenceCommand, conferenceName))
		return true
	}

	_, err = http.Post("http://redis-service:8080/saveRadiusAccountingData", jsonContentType, bytes.NewBuffer(postBody))
	if err != nil {
		log.Printf("Error sending POST request: %s\n", err)
		client.BgApi(fmt.Sprintf(destroyConferenceCommand, conferenceName))
		return true
	}

	return true
}

func handleCalleeIssue(client *goesl.Client, message *goesl.Message, operatorPrefix, destination, conferenceName string) {
	notifyCalleeIssue(message, operatorPrefix, destination)
	client.BgApi(fmt.Sprintf(destroyConferenceCommand, conferenceName))
	http.Get(fmt.Sprintf("http://redis-service:8080/popRadiusAccountingData/%s", conferenceName))
}

func notifyCalleeIssue(message *goesl.Message, operatorPrefix, destination string) {
	if message == nil {
		return
	}

	log.Printf("Callee issue: %s%q, reason - %q, code - %q",
		operatorPrefix, destination,
		message.Headers["variable_sip_invite_failure_phrase"],
		message.Headers["variable_hangup_cause_q850"])
}

func handleReadError(err error) {
	errorMessage := err.Error()
	if !strings.Contains(errorMessage, "EOF") && errorMessage != "unexpected end of JSON input" {
		log.Printf("Error reading Freeswitch message: %v", err)
	}
}

func isConnected(msg *goesl.Message, operatorPrefix, destination string) bool {
	action := msg.Headers["Action"]
	answerState := msg.Headers[answerStateHeader]
	callerDestination := msg.Headers[callerDestinationHeader]

	return (action == "add-member" && answerState == "early" && callerDestination == operatorPrefix+destination) ||
		(strings.Contains(msg.Headers["body"], operatorPrefix+destination+",conference") &&
			msg.Headers["Job-Command"] == "show" &&
			msg.Headers["Job-Command-Arg"] == "channels" &&
			msg.Headers["Event-Calling-Function"] == "bgapi_exec")
}

func isCalleeUnavailable(msg *goesl.Message, operatorPrefix, destination string) bool {
	answerState := msg.Headers[answerStateHeader]
	callerDestination := msg.Headers[callerDestinationHeader]
	hangupCauseCode := msg.Headers["variable_hangup_cause_q850"]

	return answerState == "hangup" &&
		callerDestination == operatorPrefix+destination &&
		slices.Contains(sipCalleeUnavailableCode, hangupCauseCode)
}

func isOperatorUnavailable(msg *goesl.Message, operatorPrefix, destination string) bool {
	answerState := msg.Headers[answerStateHeader]
	callerDestination := msg.Headers[callerDestinationHeader]
	hangupCauseCode := msg.Headers["variable_hangup_cause_q850"]

	return answerState == "hangup" &&
		callerDestination == operatorPrefix+destination &&
		slices.Contains(sipOperatorUnavailableCode, hangupCauseCode)
}

func logOperatorIssue(msg *goesl.Message, operatorPrefix string) {
	log.Printf(
		"%s has a problem, please contact operator %s.\n"+
			"code - %s, reason - %s",
		operatorPrefix, operatorPrefix,
		msg.Headers["variable_hangup_cause_q850"],
		msg.Headers["variable_sip_invite_failure_phrase"])
}
