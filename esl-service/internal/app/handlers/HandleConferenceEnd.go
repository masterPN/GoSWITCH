package handlers

import (
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
	hangupTimeUnix, _ := strconv.Atoi(eventData["Caller-Channel-Hangup-Time"])
	hangupTime := time.UnixMicro(int64(hangupTimeUnix))
	talkingStartTime, _ := time.Parse(timeFormat, redisData.TalkingTime)

	if nilTime, _ := time.Parse(timeFormat, "01/01/0001 00:00:00"); talkingStartTime == nilTime {
		hangupTime = nilTime
	}

	redisData.ConfID = helpers.GenerateRandomFiveDigitNumber()
	redisData.Pwd = ""
	redisData.CategoryID = "N"
	redisData.CallDuration = int(hangupTime.Sub(talkingStartTime).Seconds())
	redisData.ReleaseCode = eventData["variable_hangup_cause_q850"]
	redisData.LanguageCode = ""

	_, err := http.Post("http://mssql-service:8080/radiusAccounting", "application/json", redisData)
	if err != nil {
		log.Printf("POST http://mssql-service:8080/radiusAccounting - %s\n", err)
		return
	}
}
