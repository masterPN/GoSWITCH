package handlers

import (
	"bytes"
	"encoding/json"
	"esl-service/internal/data"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	timeFormat = "01/02/2006 15:04:05"
)

func HandleConferenceJoin(eventData map[string]string) {
	loc, _ := time.LoadLocation("Asia/Bangkok")

	callStartTimeUnix, _ := strconv.Atoi(eventData["Caller-Channel-Created-Time"])
	callAnsweredTimeUnix, _ := strconv.Atoi(eventData["Caller-Channel-Answered-Time"])
	callStartTime := time.UnixMicro(int64(callStartTimeUnix)).In(loc)
	callAnsweredTime := time.UnixMicro(int64(callAnsweredTimeUnix)).In(loc)

	radiusAccountingData := data.RadiusAccounting{
		SessionID:   eventData["variable_conference_name"],
		StartTime:   callStartTime.Format(timeFormat),
		TalkingTime: callAnsweredTime.Format(timeFormat),
	}

	saveToRedis(radiusAccountingData)
}

func saveToRedis(radiusAccountingData data.RadiusAccounting) {
	postBody, _ := json.Marshal(radiusAccountingData)

	_, err := http.Post("http://redis-service:8080/saveRadiusAccountingData", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		log.Printf("POST http://redis-service:8080/saveRadiusAccountingData - %s\n", err)
		return
	}
}