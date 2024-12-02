package onevoisdata

import (
	"database/sql"

	sharedData "github.com/masterPN/GoSWITCH-shared/data"
)

type Models struct {
	RadiusOnestageValidateData interface {
		ExecuteRadiusOnestageValidate(prefix string, callingNumber string, destinationNumber string) (sharedData.RadiusOnestageValidateData, error)
	}
	RadiusAccountingData interface {
		ExecuteRadiusAccounting(input RadiusAccounting) (RadiusAccountingData, error)
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		RadiusOnestageValidateData: RadiusOnestageValidateDataModel{DB: db},
		RadiusAccountingData:       RadiusAccountingDataModel{DB: db},
	}
}
