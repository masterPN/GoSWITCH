package onevoisdata

import (
	"database/sql"

	sharedData "github.com/masterPN/GoSWITCH-shared/data"
)

type RadiusOnestageValidateDataInterface interface {
	ExecuteRadiusOnestageValidate(prefix string, callingNumber string, destinationNumber string) (sharedData.RadiusOnestageValidateData, error)
}

type RadiusAccountingDataInterface interface {
	ExecuteRadiusAccounting(input RadiusAccounting) (RadiusAccountingData, error)
}

type Models struct {
	RadiusOnestageValidateData RadiusOnestageValidateDataInterface
	RadiusAccountingData       RadiusAccountingDataInterface
}

func NewModels(db *sql.DB) Models {
	return Models{
		RadiusOnestageValidateData: RadiusOnestageValidateDataModel{DB: db},
		RadiusAccountingData:       RadiusAccountingDataModel{DB: db},
	}
}
