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
