package handlers

import (
	"bytes"
	"encoding/json"
	"esl-service/internal/data"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	. "github.com/0x19/goesl"
)

func InitConferenceHandler(client *Client, msg map[string]string) {
	sipPort, _ := strconv.Atoi(os.Getenv("SIP_PORT"))
	externalDomain := os.Getenv("EXTERNAL_DOMAIN")
	baseClasses := strings.Split(os.Getenv("BASE_CLASS"), ",")
	operatorPrefixes := strings.Split(os.Getenv("OPERATOR_PREFIX"), ",")

	// Call is starting.
	slog.Info(msg["variable_current_application_data"])
	initConferenceData := strings.Split(msg["variable_current_application_data"], ", ")

	// Execute RadiusOnestageValidate
	postBody, _ := json.Marshal(map[string]string{
		"prefix":            "8899",
		"callingNumber":     initConferenceData[2],
		"destinationNumber": initConferenceData[3],
	})
	resp, err := http.Post("http://mssql-service:8080/radiusOnestageValidate", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		log.Printf("POST http://mssql-service:8080/radiusOnestageValidate - %s\n", err)
		return
	}
	defer resp.Body.Close()
	respBodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("POST http://mssql-service:8080/radiusOnestageValidate: could not reead response body - %s\n", err)
	}
	var radiusOnestageValidateResponse data.RadiusOnestageValidateData
	json.Unmarshal(respBodyByte, &radiusOnestageValidateResponse)

	// Break if status > 2
	// Kick A leg
	if radiusOnestageValidateResponse.Status > 2 {
		client.BgApi(fmt.Sprintf("conference %v kick all", strings.Split(initConferenceData[2], "@")[0]))
		log.Printf("status from RadiusOnestageValidate is %v\n", radiusOnestageValidateResponse.Status)
		return
	}

	// Prepare destination number, remove first 0 if contain
	if len(initConferenceData[3]) > 0 && initConferenceData[3][0] == '0' {
		// 66 stands for Thailand
		initConferenceData[3] = "66" + initConferenceData[3][1:]
	}

	// Get Operator by Number
	resp, err = http.Get(fmt.Sprintf("http://mssql-service:8080/operatorRouting?number=%s", initConferenceData[3]))
	if err != nil {
		log.Printf("GET http://mssql-service:8080/operatorRouting?number=%s - %s\n", initConferenceData[3], err)
		return
	}
	defer resp.Body.Close()
	respBodyByte, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("GET http://mssql-service:8080/operatorRouting?number=%s: could not reead response body - %s\n", initConferenceData[3], err)
	}
	var imgCdrOperatorRoutingResponse data.ImgCdrOperatorRoutingData
	json.Unmarshal(respBodyByte, &imgCdrOperatorRoutingResponse)

	baseClassResponse := [4]int{imgCdrOperatorRoutingResponse.BaseClass1, imgCdrOperatorRoutingResponse.BaseClass2, imgCdrOperatorRoutingResponse.BaseClass3, imgCdrOperatorRoutingResponse.BaseClass4}
	for _, response := range baseClassResponse {
		// skip if nil
		if response == 0 {
			continue
		}

		for j, baseClass := range baseClasses {
			if strconv.Itoa(response) == baseClass {
				// Calling B leg
				Debug("originate {origination_caller_id_number=%s}sofia/external/%s%s@%s:%v &conference(%s)",
					initConferenceData[2], operatorPrefixes[j], initConferenceData[3], externalDomain, sipPort,
					initConferenceData[2])
				client.BgApi(fmt.Sprintf("originate {origination_caller_id_number=%s}sofia/external/%s%s@%s:%v &conference(%s)",
					initConferenceData[2], operatorPrefixes[j], initConferenceData[3], externalDomain, sipPort,
					initConferenceData[2]))

				// Check B leg response within 30 seconds
				startTime := time.Now()
				for {
					// Check if 30 seconds have passed
					if time.Since(startTime) > 30*time.Second {
						break
					}

					msg, err := client.ReadMessage()

					if err != nil {

						// If it contains EOF, we really dont care...
						if !strings.Contains(err.Error(), "EOF") && err.Error() != "unexpected end of JSON input" {
							Error("Error while reading Freeswitch message: %s", err)
						}
						break
					}

					Debug("%q", msg)

					// If B receives call, then exit
					if msg.Headers["Answer-State"] == "ringing" &&
						msg.Headers["Caller-Destination-Number"] == operatorPrefixes[j]+initConferenceData[3] {
						Debug("%q receive call, then exit initConferenceHandler", operatorPrefixes[j]+initConferenceData[3])
						return
					}
				}

				return
			}
		}
	}

}
