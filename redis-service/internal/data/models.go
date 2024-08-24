package data

import "github.com/go-redis/redis/v8"

type Models struct {
	RadiusAccountingData interface {
		Set(input RadiusAccountingData) error
	}
}

func NewModels(db *redis.Client) Models {
	return Models{
		RadiusAccountingData: RadiusAccountingDataModel{DB: db},
	}
}
