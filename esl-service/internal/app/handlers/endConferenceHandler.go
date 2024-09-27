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

func EndConferenceHandler(client *goesl.Client, msg map[string]string) {
	go client.BgApi(fmt.Sprintf("conference %v kick all", msg["variable_conference_name"]))

	resp, err := http.Get(fmt.Sprintf("http://redis-service:8080/popRadiusAccountingData/%s", msg["variable_conference_name"]))
	if err != nil {
		log.Printf("GET http://redis-service:8080/popRadiusAccountingData/%s - %s\n", msg["variable_conference_name"], err)
		return
	}
	defer resp.Body.Close()
	respBodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("GET http://redis-service:8080/popRadiusAccountingData/%s: could not reead response body - %s\n", msg["variable_conference_name"], err)
	}
	var respBody data.RadiusAccountingInput
	json.Unmarshal(respBodyByte, &respBody)

	hangupTimeUnix, _ := strconv.Atoi(msg["Caller-Channel-Hangup-Time"])
	hangupTime := time.UnixMicro(int64(hangupTimeUnix))
	talkingTime, _ := time.Parse(timeFormat, respBody.TalkingTime)

	// todo: make variables dynamic
	respBody.ConfID = 65716
	respBody.Pwd = ""
	respBody.CategoryID = "N"
	respBody.CallDuration = int(hangupTime.Sub(talkingTime).Seconds())
	respBody.ReleaseCode = msg["variable_hangup_cause_q850"]
	respBody.InTrunkID = 25
	respBody.OutTrunkID = 601
	respBody.ReasonID = 0
	respBody.LanguageCode = ""

	_, err = http.Post("http://mssql-service:8080/radiusAccounting", "application/json", respBody)
	if err != nil {
		log.Printf("POST http://mssql-service:8080/radiusAccounting - %s\n", err)
		return
	}
}
