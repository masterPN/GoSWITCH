package wholesaledata

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewModels(t *testing.T) {
	mockMssqlClient := &sql.DB{}

	models := NewModels(mockMssqlClient)

	assert.NotNil(t, models.ImgCdrOperatorRoutingData)
	assert.NotNil(t, models.OptimalRouteData)
	assert.NotNil(t, models.InternalCodemappingData)

	imgCdrOperatorRoutingData, ok := models.ImgCdrOperatorRoutingData.(ImgCdrOperatorRoutingDataModel)
	assert.True(t, ok)
	assert.Equal(t, mockMssqlClient, imgCdrOperatorRoutingData.DB)

	optimalRouteData, ok := models.OptimalRouteData.(OptimalRouteDataModel)
	assert.True(t, ok)
	assert.Equal(t, mockMssqlClient, optimalRouteData.DB)

	internalCodemappingData, ok := models.InternalCodemappingData.(InternalCodemappingDataModel)
	assert.True(t, ok)
	assert.Equal(t, mockMssqlClient, internalCodemappingData.DB)
}
