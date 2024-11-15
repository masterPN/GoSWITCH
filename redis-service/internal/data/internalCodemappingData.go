package data

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type InternalCodemappingData struct {
	ID           int `json:"ID"`
	InternalCode int `json:"InternalCode"`
	OperatorCode int `json:"OperatorCode"`
}

type InternalCodemappingDataModel struct {
	DB *redis.Client
}

func (r InternalCodemappingDataModel) Set(input InternalCodemappingData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data := make(map[string]interface{})

	// Populate the data map with non-empty fields
	r.populateData(data, &input)

	// Set only the fields that are not nil or empty
	if len(data) > 0 {
		err := r.DB.HSet(ctx, "internalCodemapping:"+strconv.Itoa(input.InternalCode), data).Err()
		if err != nil {
			log.Fatalf("Could not set hash: %v", err)
			return err
		}
	}

	return nil
}

func (r InternalCodemappingDataModel) populateData(data map[string]interface{}, input *InternalCodemappingData) {
	if input.ID != 0 {
		data["ID"] = strconv.Itoa(input.ID)
	}
	if input.InternalCode != 0 {
		data["InternalCode"] = strconv.Itoa(input.InternalCode)
	}
	if input.OperatorCode != 0 {
		data["OperatorCode"] = strconv.Itoa(input.OperatorCode)
	}
}

func (r InternalCodemappingDataModel) ClearAllRecords() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	keys, err := r.DB.Keys(ctx, "internalCodemapping:*").Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		err = r.DB.Del(ctx, key).Err()
		if err != nil {
			return err
		}
	}

	return nil
}
