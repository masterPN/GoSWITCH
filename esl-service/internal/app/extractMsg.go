package app

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

const (
	eventCallingFunction = "Event-Calling-Function"
	timeFormat           = "02/01/2006 15:04:05"
)

func Execute(client *Client, msg map[string]string) {
	switch {
	case strings.Contains(msg["variable_current_application_data"], "initConference"):
		go initConferenceHandler(client, msg)
	case strings.Contains(msg["Hangup-Cause"], "CALL_REJECTED") &&
		strings.Contains(msg[eventCallingFunction], "switch_channel_perform_hangup"):
		// Callee rejects the call
		go rejectConferenceHandler(client, msg)
	case strings.Contains(msg["Answer-State"], "answered") &&
		strings.Contains(msg["Call-Direction"], "outbound") &&
		strings.Contains(msg[eventCallingFunction], "switch_channel_perform_mark_answered"):
		// Callee accepts the call
		go joinConferenceHandler(msg)
	case strings.Contains(msg["Answer-State"], "hangup") &&
		strings.Contains(msg["Call-Direction"], "inbound") &&
		strings.Contains(msg[eventCallingFunction], "switch_core_session_perform_destroy"):
		go endConferenceHandler(client, msg)
	}
}

func initConferenceHandler(client *Client, msg map[string]string) {
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

func rejectConferenceHandler(client *Client, msg map[string]string) {
	client.BgApi(fmt.Sprintf("conference %v kick all", msg["Caller-Caller-ID-Number"]))
}

func joinConferenceHandler(msg map[string]string) {
	// callerNumber := msg["Caller-Caller-ID-Number"]
	// calleeNumber := msg["Caller-Callee-ID-Number"]

	startTimeUnix, _ := strconv.Atoi(msg["Caller-Channel-Created-Time"])
	TalkingTimeUnix, _ := strconv.Atoi(msg["Caller-Channel-Answered-Time"])
	startTime := time.UnixMicro(int64(startTimeUnix))
	talkingTime := time.UnixMicro(int64(TalkingTimeUnix))

	// todo
	// make this param as dynamic, not dummy.
	postBody, _ := json.Marshal(map[string]string{
		"accessNo":     "8899",
		"anino":        "612701681",
		"destNo":       "66812424273",
		"subscriberNo": "P100000000000505",
		"sessionID":    "1723281626000100180",
		"startTime":    startTime.Format(timeFormat),
		"talkingTime":  talkingTime.Format(timeFormat),
	})

	_, err := http.Post("http://redis-service:8080/saveRadiusAccountingData", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		log.Printf("POST http://redis-service:8080/saveRadiusAccountingData - %s\n", err)
		return
	}
}

func endConferenceHandler(client *Client, msg map[string]string) {
	go client.BgApi(fmt.Sprintf("conference %v kick all", msg["variable_conference_name"]))

	// todo: anino = conference room
	anino := "612701681"

	resp, err := http.Get(fmt.Sprintf("http://redis-service:8080/popRadiusAccountingData/%s", anino))
	if err != nil {
		log.Printf("GET http://redis-service:8080/popRadiusAccountingData/%s - %s\n", anino, err)
		return
	}
	defer resp.Body.Close()
	respBodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("GET http://redis-service:8080/popRadiusAccountingData/%s: could not reead response body - %s\n", anino, err)
	}
	var respBody data.RadiusAccountingInput
	json.Unmarshal(respBodyByte, &respBody)

	// todo: make variables dynamic
	respBody.ConfID = 65716
	respBody.Pwd = ""
	respBody.CategoryID = "N"
	respBody.CallDuration, _ = strconv.Atoi(msg["variable_duration"])
	respBody.ReleaseCode = "16"
	respBody.InTrunkID = 25
	respBody.OutTrunkID = 601
	respBody.ReasonID = 0
	respBody.Prefix = ""
	respBody.LanguageCode = ""

	// todo: send post request - executeRadiusAccounting

}
