package handlers

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
)

const (
	answerStateHeader        = "Answer-State"
	callerDestinationHeader  = "Caller-Destination-Number"
	destroyConferenceCommand = "conference %v kick all"
)

var sipOperatorUnavailableCode = []string{"1", "41"}
var sipCalleeUnavailableCode = []string{"17"}

func InitConferenceHandler(client *goesl.Client, msg map[string]string) {
	sipPort, externalDomain, baseClasses, operatorPrefixes := loadConfig()

	initConferenceData := strings.Split(msg["variable_current_application_data"], ", ")

	operatorRoutingResponse, err := getOperatorRouting(prepareDestinationNumber(initConferenceData[3]))
	if err != nil {
		log.Printf("Error fetching operator routing: %s\n", err)
		return
	}

	baseClassesMap := createBaseClassesMap(baseClasses, operatorPrefixes)
	if err := originateCalls(client, initConferenceData, operatorRoutingResponse, baseClassesMap, externalDomain, sipPort, msg); err != nil {
		log.Printf("Error originating calls: %s\n", err)
	}
}

func loadConfig() (int, string, []string, []string) {
	sipPort, _ := strconv.Atoi(os.Getenv("SIP_PORT"))
	externalDomain := os.Getenv("EXTERNAL_DOMAIN")
	baseClasses := strings.Split(os.Getenv("BASE_CLASS"), ",")
	operatorPrefixes := strings.Split(os.Getenv("OPERATOR_PREFIX"), ",")
	return sipPort, externalDomain, baseClasses, operatorPrefixes
}

func validateRadius(client *goesl.Client, initConferenceData []string, msg map[string]string) bool {
	postBody, _ := json.Marshal(map[string]string{
		"prefix":            strings.Replace(msg[callerDestinationHeader], initConferenceData[3], "", 1),
		"callingNumber":     initConferenceData[2],
		"destinationNumber": initConferenceData[3],
	})
	resp, err := http.Post("http://mssql-service:8080/radiusOnestageValidate", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		log.Printf("Error validating radius: %s\n", err)
		return true
	}
	defer resp.Body.Close()

	var radiusResponse data.RadiusOnestageValidateData
	if err := json.NewDecoder(resp.Body).Decode(&radiusResponse); err != nil {
		log.Printf("Error decoding radius response: %s\n", err)
		return true
	}

	if radiusResponse.Status > 2 {
		client.BgApi(fmt.Sprintf(destroyConferenceCommand, initConferenceData[1]))
		log.Printf("Kicked due to radius status %v\n", radiusResponse.Status)
		return true
	}

	postBody, _ = json.Marshal(map[string]string{
		"accessNo":     radiusResponse.PrefixNo,
		"anino":        initConferenceData[2],
		"destNo":       initConferenceData[3],
		"subscriberNo": radiusResponse.AccountNum,
		"sessionID":    initConferenceData[1],
	})

	_, err = http.Post("http://redis-service:8080/saveRadiusAccountingData", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		log.Printf("POST http://redis-service:8080/saveRadiusAccountingData - %s\n", err)
		return true
	}

	return false
}

func prepareDestinationNumber(destination string) string {
	// Begin with 0, means local call
	if len(destination) > 0 && destination[0] == '0' {
		return "66" + destination[1:]
	}
	return destination
}

func getOperatorRouting(destination string) (data.ImgCdrOperatorRoutingData, error) {
	resp, err := http.Get(fmt.Sprintf("http://mssql-service:8080/operatorRouting?number=%s", destination))
	if err != nil {
		return data.ImgCdrOperatorRoutingData{}, err
	}
	defer resp.Body.Close()

	var routingResponse data.ImgCdrOperatorRoutingData
	if err := json.NewDecoder(resp.Body).Decode(&routingResponse); err != nil {
		return data.ImgCdrOperatorRoutingData{}, err
	}
	return routingResponse, nil
}

func createBaseClassesMap(baseClasses, operatorPrefixes []string) map[string]string {
	baseClassesMap := make(map[string]string)
	for i, v := range baseClasses {
		baseClassesMap[v] = operatorPrefixes[i]
	}
	return baseClassesMap
}

func originateCalls(client *goesl.Client, initConferenceData []string, routingResponse data.ImgCdrOperatorRoutingData, baseClassesMap map[string]string, externalDomain string, sipPort int, msg map[string]string) error {
	baseClassResponse := [4]int{
		routingResponse.BaseClass1,
		routingResponse.BaseClass2,
		routingResponse.BaseClass3,
		routingResponse.BaseClass4,
	}

	for i, response := range baseClassResponse {
		if response == 0 {
			continue
		}

		if operatorPrefix, exists := baseClassesMap[strconv.Itoa(response)]; exists && operatorPrefix != "" {
			if validateRadius(client, initConferenceData, msg) {
				continue
			}

			if originateCall(client, initConferenceData, operatorPrefix, externalDomain, sipPort) {
				return nil
			}
		}

		if i == len(baseClassResponse)-1 {
			goesl.Debug("There's no operator available.")
			client.BgApi(fmt.Sprintf(destroyConferenceCommand, initConferenceData[1]))
			http.Get(fmt.Sprintf("http://redis-service:8080/popRadiusAccountingData/%s", initConferenceData[1]))
		}
	}
	return nil
}

func originateCall(client *goesl.Client, initConferenceData []string, operatorPrefix, externalDomain string, sipPort int) bool {
	client.BgApi(fmt.Sprintf("originate {origination_caller_id_number=%s}sofia/external/%s%s@%s:%v &conference(%s)",
		initConferenceData[2], operatorPrefix, prepareDestinationNumber(initConferenceData[3]), externalDomain, sipPort,
		initConferenceData[1]))

	return waitForCall(client, operatorPrefix, prepareDestinationNumber(initConferenceData[3]), initConferenceData[1])
}

func waitForCall(client *goesl.Client, operatorPrefix, destination string, conferenceName string) bool {
	startTime := time.Now()
	var once sync.Once

	for {
		if time.Since(startTime) > 5*time.Second {
			break
		} else if time.Since(startTime) > 2*time.Second {
			go once.Do(func() {
				client.BgApi("show channels")
			})
		}

		msg, err := client.ReadMessage()
		if err != nil {
			handleReadError(err)
			break
		}

		if isConnected(msg, operatorPrefix, destination) {
			goesl.Debug("%q received call, exiting initConferenceHandler", operatorPrefix+destination)
			return true
		}

		if isCalleeUnavailable(msg, operatorPrefix, destination) {
			logCalleeIssue(msg, operatorPrefix, destination)
			client.BgApi(fmt.Sprintf(destroyConferenceCommand, conferenceName))
			http.Get(fmt.Sprintf("http://redis-service:8080/popRadiusAccountingData/%s", conferenceName))
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
