package data

import (
	"context"
	"fmt"
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

const internalCodemappingKeyPrefix = "internalCodemapping:"

func (r InternalCodemappingDataModel) Set(input InternalCodemappingData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data := make(map[string]interface{})

	// Populate the data map with non-empty fields
	r.populateData(data, &input)

	// Set only the fields that are not nil or empty
	if len(data) > 0 {
		err := r.DB.HSet(ctx, internalCodemappingKeyPrefix+strconv.Itoa(input.InternalCode), data).Err()
		if err != nil {
			log.Printf("Error setting hash: %v", err)
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

func (r InternalCodemappingDataModel) Get(internalCode int) (InternalCodemappingData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	internalCodeMappingData, err := r.DB.HGetAll(ctx, internalCodemappingKeyPrefix+strconv.Itoa(internalCode)).Result()
	if err != nil {
		return InternalCodemappingData{}, fmt.Errorf("failed to get internal code mapping: %w", err)
	}

	if len(internalCodeMappingData) == 0 {
		return InternalCodemappingData{}, fmt.Errorf("internal code mapping not found")
	}

	id, _ := strconv.Atoi(internalCodeMappingData["ID"])
	operatorCode, _ := strconv.Atoi(internalCodeMappingData["OperatorCode"])

	result := InternalCodemappingData{
		ID:           id,
		InternalCode: internalCode,
		OperatorCode: operatorCode,
	}

	return result, nil
}

func (r InternalCodemappingDataModel) Delete(internalCode int) error {
	if internalCode <= 0 {
		return fmt.Errorf("invalid internal code: %d", internalCode)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.DB.Del(ctx, internalCodemappingKeyPrefix+strconv.Itoa(internalCode)).Err()
	if err != nil {
		return fmt.Errorf("failed to delete internal code mapping: %w", err)
	}

	return nil
}
