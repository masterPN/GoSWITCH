package wholesaledata

import "database/sql"

type ImgCdrOperatorRoutingDataInterface interface {
	GetFirstImgCdrOperatorRoutingByNumber(number string) (ImgCdrOperatorRoutingData, error)
}

type OptimalRouteDataInterface interface {
	ExecuteGetOptimalRoute(pCallString string) (OptimalRouteData, error)
}

type InternalCodemappingDataInterface interface {
	GetAll() ([]InternalCodemappingData, error)
	Set(input InternalCodemappingData) (InternalCodemappingData, error)
	Delete(internalCode int) error
}

type Models struct {
	ImgCdrOperatorRoutingData ImgCdrOperatorRoutingDataInterface
	OptimalRouteData          OptimalRouteDataInterface
	InternalCodemappingData   InternalCodemappingDataInterface
}

func NewModels(db *sql.DB) Models {
	return Models{
		ImgCdrOperatorRoutingData: ImgCdrOperatorRoutingDataModel{DB: db},
		OptimalRouteData:          OptimalRouteDataModel{DB: db},
		InternalCodemappingData:   InternalCodemappingDataModel{DB: db},
	}
}
