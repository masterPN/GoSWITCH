package handlers

import (
	"encoding/json"
	"esl-service/internal/data"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	. "github.com/0x19/goesl"
)

func EndConferenceHandler(client *Client, msg map[string]string) {
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

	resp, err = http.Post("http://mssql-service:8080/radiusAccounting", "application/json", respBody)
	if err != nil {
		log.Printf("POST http://mssql-service:8080/radiusAccounting - %s\n", err)
		return
	}
}
