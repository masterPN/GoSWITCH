package handlers

import (
	"bytes"
	"encoding/json"
	"esl-service/internal/app/helpers"
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

	redisData := getRedisData(eventData["variable_conference_name"])
	if redisData == nil {
		return
	}

	updateMsSql(redisData, eventData)
}

func getRedisData(conferenceName string) *data.RadiusAccounting {
	redisResponse, err := http.Get(fmt.Sprintf("http://redis-service:8080/popRadiusAccountingData/%s", conferenceName))
	if err != nil {
		log.Printf("GET http://redis-service:8080/popRadiusAccountingData/%s - %s\n", conferenceName, err)
		return nil
	}
	defer redisResponse.Body.Close()

	redisResponseBodyBytes, err := io.ReadAll(redisResponse.Body)
	if err != nil {
		log.Printf("GET http://redis-service:8080/popRadiusAccountingData/%s: could not read response body - %s\n", conferenceName, err)
		return nil
	}

	var redisResponseData data.RadiusAccounting
	json.Unmarshal(redisResponseBodyBytes, &redisResponseData)

	return &redisResponseData
}

func updateMsSql(redisData *data.RadiusAccounting, eventData map[string]string) {
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Printf("Error loading location: %s\n", err)
		return
	}

	hangupTimeUnix, err := strconv.Atoi(eventData["Caller-Channel-Hangup-Time"])
	if err != nil {
		log.Printf("Error converting hangup time: %s\n", err)
		return
	}
	hangupTime := time.UnixMicro(int64(hangupTimeUnix)).In(location)

	talkingStartTime, err := time.ParseInLocation(timeFormat, redisData.TalkingTime, location)
	if err != nil {
		log.Printf("Error parsing talking start time: %s\n", err)
		return
	}

	nilTime, _ := time.Parse(timeFormat, "01/01/0001 00:00:00")
	if talkingStartTime.Equal(nilTime) {
		hangupTime = nilTime
	}

	redisData.ConfID = helpers.GenerateRandomFiveDigitNumber()
	redisData.Pwd = ""
	redisData.CategoryID = "N"
	redisData.CallDuration = int(hangupTime.Sub(talkingStartTime).Seconds())
	redisData.ReleaseCode = eventData["variable_hangup_cause_q850"]
	redisData.LanguageCode = ""

	jsonData, err := json.Marshal(redisData)
	if err != nil {
		log.Printf("Error marshaling redisData to JSON: %s\n", err)
		return
	}

	response, err := http.Post("http://mssql-service:8080/radiusAccounting", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending POST request: %s\n", err)
		return
	}
	defer response.Body.Close()
}
