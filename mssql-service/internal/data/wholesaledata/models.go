package wholesaledata

import "database/sql"

type Models struct {
	ImgCdrOperatorRoutingData interface {
		GetFirstImgCdrOperatorRoutingByNumber(number string) (ImgCdrOperatorRoutingData, error)
	}
	OptimalRouteData interface {
		ExecuteGetOptimalRoute(pCallString string) (OptimalRouteData, error)
	}
	InternalCodemappingData interface {
		GetAll() ([]InternalCodemappingData, error)
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		ImgCdrOperatorRoutingData: ImgCdrOperatorRoutingDataModel{DB: db},
		OptimalRouteData:          OptimalRouteDataModel{DB: db},
		InternalCodemappingData:   InternalCodemappingDataModel{DB: db},
	}
}
