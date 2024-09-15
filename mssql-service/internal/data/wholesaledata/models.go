package wholesaledata

import "database/sql"

type Models struct {
	ImgCdrOperatorRoutingData interface {
		GetFirstImgCdrOperatorRoutingByNumber(number string) (ImgCdrOperatorRoutingData, error)
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		ImgCdrOperatorRoutingData: ImgCdrOperatorRoutingDataModel{DB: db},
	}
}
