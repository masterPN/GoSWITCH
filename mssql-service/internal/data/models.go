package data

import "database/sql"

type Models struct {
	RadiusData interface {
		ExecuteRadiusOnestageValidate(prefix string, callingNumber string, destinationNumber string) (RadiusData, error)
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		RadiusData: RadiusDataModel{DB: db},
	}
}
