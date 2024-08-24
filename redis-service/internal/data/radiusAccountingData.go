package data

import (
	"context"
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
