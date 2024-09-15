package onevoisdata

import "database/sql"

type Models struct {
	RadiusOnestageValidateData interface {
		ExecuteRadiusOnestageValidate(prefix string, callingNumber string, destinationNumber string) (RadiusOnestageValidateData, error)
	}
	RadiusAccountingData interface {
		ExecuteRadiusAccounting(input RadiusAccountingInput) (RadiusAccountingData, error)
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		RadiusOnestageValidateData: RadiusOnestageValidateDataModel{DB: db},
		RadiusAccountingData:       RadiusAccountingDataModel{DB: db},
	}
}
