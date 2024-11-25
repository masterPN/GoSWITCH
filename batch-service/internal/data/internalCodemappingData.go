package data

import (
	"batch-service/internal/helpers"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type InternalCodemappingData struct {
	ID           int `json:"ID"`
	InternalCode int `json:"InternalCode"`
	OperatorCode int `json:"OperatorCode"`
}

// Read implements io.Reader.
func (i InternalCodemappingData) Read(p []byte) (n int, err error) {
	// Convert the struct fields to a byte slice
	data, err := json.Marshal(i)
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

func (i InternalCodemappingData) SendInternalCodemappingDataToRedis() error {
	url := "http://redis-service:8080/internalCodemappingData"
	resp, err := helpers.PostRequest(url, i)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP request failed: %s", string(bodyBytes))
	}

	return nil
}

func (i InternalCodemappingData) SendInternalCodemappingDataToMssql() error {
	url := "http://mssql-service:8080/internalCodemapping"
	resp, err := helpers.PostRequest(url, i)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP request failed: %s", string(bodyBytes))
	}

	return nil
}
