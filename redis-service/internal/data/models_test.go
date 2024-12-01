package data

import (
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestNewModels(t *testing.T) {
	// Create a mock Redis client using redis
	mockRedisClient := &redis.Client{}

	// Call NewModels with the mock Redis client
	models := NewModels(mockRedisClient)

	// Assert that the RadiusAccountingData and InternalCodemappingData fields are not nil
	assert.NotNil(t, models.RadiusAccountingData)
	assert.NotNil(t, models.InternalCodemappingData)

	// Assert that the RadiusAccountingData and InternalCodemappingData fields are instances of the respective models
	// Assert that the Redis client passed to the models is the same as the mock Redis client
	radiusAccountingData, ok := models.RadiusAccountingData.(RadiusAccountingDataModel)
	assert.True(t, ok)
	assert.Equal(t, mockRedisClient, radiusAccountingData.DB)

	internalCodemappingData, ok := models.InternalCodemappingData.(InternalCodemappingDataModel)
	assert.True(t, ok)
	assert.Equal(t, mockRedisClient, internalCodemappingData.DB)
}
