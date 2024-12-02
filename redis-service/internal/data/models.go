package data

import (
	"github.com/go-redis/redis/v8"
)

// Define the interfaces as named types
type RadiusAccountingDataInterface interface {
	Set(input RadiusAccountingData) error
	Pop(anino string) (RadiusAccountingData, error)
}

type InternalCodemappingDataInterface interface {
	Set(input InternalCodemappingData) error
	Get(internalCode int) (InternalCodemappingData, error)
	Delete(internalCode int) error
}

// Now use these interfaces in your Models struct
type Models struct {
	RadiusAccountingData    RadiusAccountingDataInterface
	InternalCodemappingData InternalCodemappingDataInterface
}

func NewModels(db *redis.Client) Models {
	return Models{
		RadiusAccountingData:    RadiusAccountingDataModel{DB: db},
		InternalCodemappingData: InternalCodemappingDataModel{DB: db},
	}
}
