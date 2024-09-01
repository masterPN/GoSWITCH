package data

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RadiusAccountingData struct {
	AccessNo     string `json:"accessNo"`
	Anino        string `json:"anino"`
	DestNo       string `json:"destNo"`
	SubscriberNo string `json:"subscriberNo"`
	SessionID    string `json:"sessionID"`
	StartTime    string `json:"startTime"`
	TalkingTime  string `json:"talkingTime"`
}

type RadiusAccountingDataModel struct {
	DB *redis.Client
}

func (r RadiusAccountingDataModel) Set(input RadiusAccountingData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.DB.HSet(ctx, input.Anino, map[string]interface{}{
		"accessNo":     input.AccessNo,
		"anino":        input.Anino,
		"destNo":       input.DestNo,
		"subscriberNo": input.SubscriberNo,
		"sessionID":    input.SessionID,
		"startTime":    input.StartTime,
		"talkingTime":  input.TalkingTime,
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

	radiusAccountingData := RadiusAccountingData{
		AccessNo:     radiusAccountingDataMap["accessNo"],
		Anino:        radiusAccountingDataMap["anino"],
		DestNo:       radiusAccountingDataMap["destNo"],
		SubscriberNo: radiusAccountingDataMap["subscriberNo"],
		SessionID:    radiusAccountingDataMap["sessionID"],
		StartTime:    radiusAccountingDataMap["startTime"],
		TalkingTime:  radiusAccountingDataMap["talkingTime"],
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
		return radiusAccountingData, fmt.Errorf("Key %s does not exist.\n", anino)
	}

	return radiusAccountingData, nil
}
