package handlers

import (
	"bytes"
	"encoding/json"
	"esl-service/internal/data"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/0x19/goesl"
)

func InitConferenceHandler(client *goesl.Client, msg map[string]string) {
	sipPort, externalDomain, baseClasses, operatorPrefixes := loadConfig()

	initConferenceData := strings.Split(msg["variable_current_application_data"], ", ")
	if validateRadius(client, initConferenceData) {
		return
	}

	// Prepare destination number
	initConferenceData[3] = prepareDestinationNumber(initConferenceData[3])

	operatorRoutingResponse, err := getOperatorRouting(initConferenceData[3])
	if err != nil {
		log.Printf("Error fetching operator routing: %s\n", err)
		return
	}

	baseClassesMap := createBaseClassesMap(baseClasses, operatorPrefixes)
	if err := originateCalls(client, initConferenceData, operatorRoutingResponse, baseClassesMap, externalDomain, sipPort); err != nil {
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

func validateRadius(client *goesl.Client, initConferenceData []string) bool {
	postBody, _ := json.Marshal(map[string]string{
		"prefix":            "8899",
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
		client.BgApi(fmt.Sprintf("conference %v kick all", strings.Split(initConferenceData[2], "@")[0]))
		log.Printf("Kicked due to radius status %v\n", radiusResponse.Status)
		return true
	}
	return false
}

func prepareDestinationNumber(destination string) string {
	if len(destination) > 0 && destination[0] == '0' {
		return destination[1:]
	}
	return destination
}

func getOperatorRouting(destination string) (data.ImgCdrOperatorRoutingData, error) {
	resp, err := http.Get(fmt.Sprintf("http://mssql-service:8080/operatorRouting?number=%s", "66"+destination))
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

func originateCalls(client *goesl.Client, initConferenceData []string, routingResponse data.ImgCdrOperatorRoutingData, baseClassesMap map[string]string, externalDomain string, sipPort int) error {
	baseClassResponse := [4]int{
		routingResponse.BaseClass1,
		routingResponse.BaseClass2,
		routingResponse.BaseClass3,
		routingResponse.BaseClass4,
	}

	for _, response := range baseClassResponse {
		if response == 0 {
			continue
		}

		if operatorPrefix, exists := baseClassesMap[strconv.Itoa(response)]; exists && operatorPrefix != "" {
			if err := originateCall(client, initConferenceData, operatorPrefix, externalDomain, sipPort); err != nil {
				return err
			}
		}
	}
	return nil
}

func originateCall(client *goesl.Client, initConferenceData []string, operatorPrefix, externalDomain string, sipPort int) error {
	client.BgApi(fmt.Sprintf("originate {origination_caller_id_number=%s}sofia/external/%s%s@%s:%v &conference(%s)",
		initConferenceData[2], operatorPrefix, initConferenceData[3], externalDomain, sipPort,
		initConferenceData[2]))

	return waitForCall(client, operatorPrefix, initConferenceData[3])
}

func waitForCall(client *goesl.Client, operatorPrefix, destination string) error {
	startTime := time.Now()
	for {
		if time.Since(startTime) > 5*time.Second {
			break
		}

		msg, err := client.ReadMessage()
		if err != nil {
			if !strings.Contains(err.Error(), "EOF") && err.Error() != "unexpected end of JSON input" {
				goesl.Error("Error while reading Freeswitch message: %s", err)
			}
			break
		}

		if msg.Headers["Action"] == "add-member" && msg.Headers["Answer-State"] == "early" && msg.Headers["Caller-Destination-Number"] == operatorPrefix+destination {
			goesl.Debug("%q received call, exiting initConferenceHandler", operatorPrefix+destination)
			return nil
		}
	}
	return nil
}
