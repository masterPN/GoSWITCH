package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func MakeRequest(url string, method string, contentType string, data interface{}) (*http.Response, error) {
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func PostRequest(url string, i interface{}) (*http.Response, error) {
	resp, err := MakeRequest(url, http.MethodPost, "application/json", i)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP request failed: %s", string(bodyBytes))
	}

	return resp, nil
}

func DeleteRequest(url string, i interface{}) (*http.Response, error) {
	resp, err := MakeRequest(url, http.MethodDelete, "application/json", i)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP request failed: %s", string(bodyBytes))
	}

	return resp, nil
}
