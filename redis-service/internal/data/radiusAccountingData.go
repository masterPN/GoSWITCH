package data

import (
	"context"
	"fmt"
	"log"
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

	err := r.DB.HSet(ctx, input.SessionID, map[string]interface{}{
		"confID":       input.ConfID,
		"accessNo":     input.AccessNo,
		"anino":        input.Anino,
		"destNo":       input.DestNo,
		"subscriberNo": input.SubscriberNo,
		"pwd":          input.Pwd,
		"sessionID":    input.SessionID,
		"categoryID":   input.CategoryID,
		"startTime":    input.StartTime,
		"talkingTime":  input.TalkingTime,
		"callDuration": input.CallDuration,
		"releaseCode":  input.ReleaseCode,
		"inTrunkID":    input.InTrunkID,
		"outTrunkID":   input.OutTrunkID,
		"reasonID":     input.ReasonID,
		"prefix":       input.Prefix,
		"languageCode": input.LanguageCode,
	}).Err()
	if err != nil {
		log.Fatalf("Could not set hash: %v", err)
		return err
	}

	return nil
}

func (r RadiusAccountingDataModel) Pop(anino string) (RadiusAccountingData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Retrieve all fields and values from the hash
	radiusAccountingDataMap, err := r.DB.HGetAll(ctx, anino).Result()
	if err != nil {
		log.Fatalf("could not HGetAll from hash %s: %v", anino, err)
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
	result, err := r.DB.Del(ctx, anino).Result()

	if err != nil {
		log.Fatalf("could not delete key %s: %v", anino, err)
		return radiusAccountingData, err
	}

	// Output the result
	if result == 1 {
		fmt.Printf("Key %s was deleted successfully.\n", anino)
	} else {
		fmt.Printf("Key %s does not exist.\n", anino)
		return radiusAccountingData, fmt.Errorf("key %s does not exist", anino)
	}

	return radiusAccountingData, nil
}
