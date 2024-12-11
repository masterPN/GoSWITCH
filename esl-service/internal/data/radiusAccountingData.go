package data

import (
	"encoding/json"
	"io"
)

type RadiusAccounting struct {
	ConfID       int    `json:"confID"`
	AccessNo     string `json:"accessNo"`
	Anino        string `json:"anino"`
	DestNo       string `json:"destNo"`
	SubscriberNo string `json:"subscriberNo"`
	Pwd          string `json:"pwd"`
	SessionID    string `json:"sessionID"`
	CategoryID   string `json:"categoryID"`
	StartTime    string `json:"startTime"`
	TalkingTime  string `json:"talkingTime"`
	CallDuration int    `json:"callDuration"`
	ReleaseCode  string `json:"releaseCode"`
	InTrunkID    int    `json:"inTrunkID"`
	OutTrunkID   int    `json:"outTrunkID"`
	ReasonID     int    `json:"reasonID"`
	Prefix       string `json:"prefix"`
	LanguageCode string `json:"languageCode"`
}

// Read implements io.Reader.
func (r RadiusAccounting) Read(p []byte) (n int, err error) {
	// Serialize the struct to JSON.
	data, err := json.Marshal(r)
	if err != nil {
		return 0, err
	}

	// Copy the serialized data into the provided byte slice p.
	// The number of bytes to copy is the minimum of the available space in p and the length of the JSON data.
	copyLen := len(data)
	if copyLen > len(p) {
		copyLen = len(p)
	}
	copy(p, data[:copyLen])

	// Return the number of bytes copied and an error if the data was truncated.
	if copyLen < len(data) {
		return copyLen, io.ErrShortBuffer
	}
	return copyLen, nil
}
