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

func JoinConferenceHandler(msg map[string]string) {
	startTimeUnix, _ := strconv.Atoi(msg["Caller-Channel-Created-Time"])
	TalkingTimeUnix, _ := strconv.Atoi(msg["Caller-Channel-Answered-Time"])
	startTime := time.UnixMicro(int64(startTimeUnix))
	talkingTime := time.UnixMicro(int64(TalkingTimeUnix))

	radiusAccountingBody := data.RadiusAccounting{
		SessionID:   msg["variable_conference_name"],
		StartTime:   startTime.Format(timeFormat),
		TalkingTime: talkingTime.Format(timeFormat),
	}

	postBody, _ := json.Marshal(radiusAccountingBody)

	_, err := http.Post("http://redis-service:8080/saveRadiusAccountingData", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		log.Printf("POST http://redis-service:8080/saveRadiusAccountingData - %s\n", err)
		return
	}
}
