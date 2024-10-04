package handlers

import (
	"encoding/json"
	"esl-service/internal/data"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/0x19/goesl"
)

func HandleConferenceEnd(eslClient *goesl.Client, eventData map[string]string) {
	go eslClient.BgApi(fmt.Sprintf("conference %v kick all", eventData["variable_conference_name"]))

	redisResponse, err := http.Get(fmt.Sprintf("http://redis-service:8080/popRadiusAccountingData/%s", eventData["variable_conference_name"]))
	if err != nil {
		log.Printf("GET http://redis-service:8080/popRadiusAccountingData/%s - %s\n", eventData["variable_conference_name"], err)
		return
	}
	defer redisResponse.Body.Close()
	redisResponseBodyBytes, err := io.ReadAll(redisResponse.Body)
	if err != nil {
		log.Printf("GET http://redis-service:8080/popRadiusAccountingData/%s: could not read response body - %s\n", eventData["variable_conference_name"], err)
	}
	var redisResponseData data.RadiusAccounting
	json.Unmarshal(redisResponseBodyBytes, &redisResponseData)

	hangupTimeUnix, _ := strconv.Atoi(eventData["Caller-Channel-Hangup-Time"])
	hangupTime := time.UnixMicro(int64(hangupTimeUnix))
	talkingStartTime, _ := time.Parse(timeFormat, redisResponseData.TalkingTime)

	if nilTime, _ := time.Parse(timeFormat, "01/01/0001 00:00:00"); talkingStartTime == nilTime {
		hangupTime = nilTime
	}

	// todo: make variables dynamic
	redisResponseData.ConfID = 65716
	redisResponseData.Pwd = ""
	redisResponseData.CategoryID = "N"
	redisResponseData.CallDuration = int(hangupTime.Sub(talkingStartTime).Seconds())
	redisResponseData.ReleaseCode = eventData["variable_hangup_cause_q850"]
	redisResponseData.LanguageCode = ""

	_, err = http.Post("http://mssql-service:8080/radiusAccounting", "application/json", redisResponseData)
	if err != nil {
		log.Printf("POST http://mssql-service:8080/radiusAccounting - %s\n", err)
		return
	}
}
