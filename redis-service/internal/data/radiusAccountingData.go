package data

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RadiusAccountingData struct {
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

type RadiusAccountingDataModel struct {
	DB *redis.Client
}

func (r RadiusAccountingDataModel) Set(input RadiusAccountingData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data := make(map[string]interface{})

	// Populate the data map with non-empty fields
	r.populateData(data, input)

	// Set only the fields that are not nil or empty
	if len(data) > 0 {
		err := r.DB.HSet(ctx, input.SessionID, data).Err()
		if err != nil {
			log.Printf("Could not set hash: %v", err)
			return err
		}
	}

	return nil
}

// populateData fills the provided map with non-empty fields from the input.
func (r RadiusAccountingDataModel) populateData(data map[string]interface{}, input RadiusAccountingData) {
	value := reflect.ValueOf(input)
	typeOf := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldName := typeOf.Field(i).Name

		switch field.Kind() {
		case reflect.Int:
			if field.Int() != 0 {
				data[fieldName] = field.Int()
			}
		case reflect.String:
			if field.String() != "" {
				data[fieldName] = field.String()
			}
		}
	}
}

func (r RadiusAccountingDataModel) Pop(sessionID string) (RadiusAccountingData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Retrieve all fields and values from the hash
	radiusAccountingDataMap, err := r.DB.HGetAll(ctx, sessionID).Result()
	if err != nil {
		log.Printf("could not HGetAll from hash %s: %v", sessionID, err)
		return RadiusAccountingData{}, err
	}

	// Convert string fields to appropriate types
	confID, _ := strconv.Atoi(radiusAccountingDataMap["confID"])
	callDuration, _ := strconv.Atoi(radiusAccountingDataMap["callDuration"])
	inTrunkID, _ := strconv.Atoi(radiusAccountingDataMap["inTrunkID"])
	outTrunkID, _ := strconv.Atoi(radiusAccountingDataMap["outTrunkID"])
	reasonID, _ := strconv.Atoi(radiusAccountingDataMap["reasonID"])

	radiusAccountingData := RadiusAccountingData{
		ConfID:       confID,
		AccessNo:     radiusAccountingDataMap["accessNo"],
		Anino:        radiusAccountingDataMap["anino"],
		DestNo:       radiusAccountingDataMap["destNo"],
		SubscriberNo: radiusAccountingDataMap["subscriberNo"],
		Pwd:          radiusAccountingDataMap["pwd"],
		SessionID:    radiusAccountingDataMap["sessionID"],
		CategoryID:   radiusAccountingDataMap["categoryID"],
		StartTime:    radiusAccountingDataMap["startTime"],
		TalkingTime:  radiusAccountingDataMap["talkingTime"],
		CallDuration: callDuration,
		ReleaseCode:  radiusAccountingDataMap["releaseCode"],
		InTrunkID:    inTrunkID,
		OutTrunkID:   outTrunkID,
		ReasonID:     reasonID,
		Prefix:       radiusAccountingDataMap["prefix"],
		LanguageCode: radiusAccountingDataMap["languageCode"],
	}

	// Delete the key
	result, err := r.DB.Del(ctx, sessionID).Result()

	if err != nil {
		log.Printf("could not delete key %s: %v", sessionID, err)
		return radiusAccountingData, err
	}

	// Output the result
	if result == 1 {
		fmt.Printf("Key %s was deleted successfully.\n", sessionID)
	} else {
		fmt.Printf("Key %s does not exist.\n", sessionID)
		return radiusAccountingData, fmt.Errorf("key %s does not exist", sessionID)
	}

	return radiusAccountingData, nil
}
