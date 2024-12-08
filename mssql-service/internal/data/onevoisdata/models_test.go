package onevoisdata

import (
	"testing"

	"database/sql"

	"github.com/stretchr/testify/assert"
)

func TestNewModels(t *testing.T) {
	mockMssqlClient := &sql.DB{}

	models := NewModels(mockMssqlClient)

	assert.NotNil(t, models.RadiusOnestageValidateData)
	assert.NotNil(t, models.RadiusAccountingData)

	radiusOnestageValidateData, ok := models.RadiusOnestageValidateData.(RadiusOnestageValidateDataModel)
	assert.True(t, ok)
	assert.Equal(t, mockMssqlClient, radiusOnestageValidateData.DB)

	radiusAccountingData, ok := models.RadiusAccountingData.(RadiusAccountingDataModel)
	assert.True(t, ok)
	assert.Equal(t, mockMssqlClient, radiusAccountingData.DB)
}
