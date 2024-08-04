package data

import "database/sql"

type Models struct {
	RadiusData interface {
		ExecuteRadiusOnestageValidate(prefix int, callingNumber int, destinationNumber int) (RadiusData, error)
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		RadiusData: RadiusDataModel{DB: db},
	}
}
