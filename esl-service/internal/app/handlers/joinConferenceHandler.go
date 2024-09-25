package handlers

import (
	"bytes"
	"encoding/json"
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

	postBody, _ := json.Marshal(map[string]string{
		"sessionID":   msg["variable_conference_name"],
		"startTime":   startTime.Format(timeFormat),
		"talkingTime": talkingTime.Format(timeFormat),
	})

	_, err := http.Post("http://redis-service:8080/saveRadiusAccountingData", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		log.Printf("POST http://redis-service:8080/saveRadiusAccountingData - %s\n", err)
		return
	}
}
