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

	. "github.com/0x19/goesl"
)

func InitConferenceHandler(client *Client, msg map[string]string) {
	sipPort, _ := strconv.Atoi(os.Getenv("SIP_PORT"))

	// Call is starting.
	slog.Info(msg["variable_current_application_data"])
	initConferenceData := strings.Split(msg["variable_current_application_data"], ", ")

	// Execute RadiusOnestageValidate
	// postBody, _ := json.Marshal(map[string]string{
	// 	"prefix":            "8899",
	// 	"callingNumber":     initConferenceData[2],
	// 	"destinationNumber": initConferenceData[3],
	// })

	// todo
	// make this param as dynamic, not dummy.
	postBody, _ := json.Marshal(map[string]string{
		"prefix":            "8899",
		"callingNumber":     "6627288000",
		"destinationNumber": "0844385742",
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
	var respBody data.RadiusData
	json.Unmarshal(respBodyByte, &respBody)

	// Break if status != 0
	// Kick A leg
	if respBody.Status != 0 {
		client.BgApi(fmt.Sprintf("conference %v kick all", strings.Split(initConferenceData[2], "@")[0]))
		log.Printf("status from RadiusOnestageValidate is %v\n", respBody.Status)
		return
	}

	// Calling B leg
	client.BgApi(fmt.Sprintf("originate {origination_caller_id_number=%s}sofia/internal/%s:%v &conference(%s)",
		initConferenceData[2], initConferenceData[3], sipPort,
		initConferenceData[2]))
}
