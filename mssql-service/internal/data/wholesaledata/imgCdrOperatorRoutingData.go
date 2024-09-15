package wholesaledata

import (
	"context"
	"database/sql"
	"time"
)

type ImgCdrOperatorRoutingData struct {
	RoutingPrefix              string     `json:"RoutingPrefix"`
	BaseOperator               *string    `json:"BaseOperator"`
	BaseClass1                 *int       `json:"BaseClass1"`
	BaseClass2                 *int       `json:"BaseClass2"`
	BaseClass3                 *int       `json:"BaseClass3"`
	BaseClass4                 *int       `json:"BaseClass4"`
	BaseClassCreationDate      *time.Time `json:"BaseClassCreationDate"`
	CurrentClass1              *int       `json:"CurrentClass1"`
	CurrentClass2              *int       `json:"CurrentClass2"`
	CurrentClass3              *int       `json:"CurrentClass3"`
	CurrentClass4              *int       `json:"CurrentClass4"`
	CurrentClassLastUpdateDate *time.Time `json:"CurrentClassLastUpdateDate"`
	CurrentClassRemarks        *string    `json:"CurrentClassRemarks"`
	RoutingManualOverride      *string    `json:"RoutingManualOverride"`
	OverrideCreationDate       *time.Time `json:"OverrideCreationDate"`
	OverrideCreatedBy          *string    `json:"OverrideCreatedBy"`
	OverrideReason             *string    `json:"OverrideReason"`
	BaseClassLastUpdateDate    *time.Time `json:"BaseClassLastUpdateDate"`
}

type ImgCdrOperatorRoutingDataModel struct {
	DB *sql.DB
}

func (r ImgCdrOperatorRoutingDataModel) GetFirstImgCdrOperatorRoutingByNumber(number string) (ImgCdrOperatorRoutingData, error) {
	query := `SELECT
					TOP 1 *
				FROM
					IMGCDROPERATORROUTING IR
				WHERE
					$1 LIKE CONCAT(IR.ROUTINGPREFIX, '%')
				ORDER BY
					IR.ROUTINGPREFIX DESC;`

	args := []interface{}{
		number,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := r.DB.QueryRowContext(ctx, query, args...)

	var result ImgCdrOperatorRoutingData
	err := row.Scan(&result.RoutingPrefix, &result.BaseOperator, &result.BaseClass1, &result.BaseClass2, &result.BaseClass3, &result.BaseClass4, &result.BaseClassCreationDate, &result.CurrentClass1, &result.CurrentClass2, &result.CurrentClass3, &result.CurrentClass4, &result.CurrentClassLastUpdateDate, &result.CurrentClassRemarks, &result.RoutingManualOverride, &result.OverrideCreationDate, &result.OverrideCreatedBy, &result.OverrideReason, &result.BaseClassLastUpdateDate)
	if err != nil {
		return ImgCdrOperatorRoutingData{}, err
	}

	return result, nil
}
