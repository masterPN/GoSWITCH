package onevoisdata

import (
	"context"
	"database/sql"
	"time"

	sharedData "github.com/masterPN/GoSWITCH-shared/data"
)

type RadiusOnestageValidateDataModel struct {
	DB *sql.DB
}

func (r RadiusOnestageValidateDataModel) ExecuteRadiusOnestageValidate(prefix string, callingNumber string, destinationNumber string) (sharedData.RadiusOnestageValidateData, error) {
	query := `EXEC RADIUS_ONESTAGE_VALIDATE @SESSION_ID = $1,
											@ACCESS_NO 	= $2,
											@ANINO 		= $3,
											@DIDNO 		= $4;`

	args := []interface{}{
		prefix + destinationNumber,
		prefix,
		callingNumber,
		destinationNumber,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := r.DB.QueryRowContext(ctx, query, args...)

	var result sharedData.RadiusOnestageValidateData
	err := row.Scan(&result.Lcode, &result.Status, &result.RouteType, &result.CallType, &result.AccountNum, &result.PrefixNo, &result.Dnis, &result.Pin, &result.FollowOnCall, &result.Trunk1, &result.Carrier1, &result.PlanCode1, &result.Plan1, &result.Duration1, &result.Trunk2, &result.Carrier2, &result.PlanCode2, &result.Plan2, &result.Duration2, &result.Trunk3, &result.Carrier3, &result.PlanCode3, &result.Plan3, &result.Duration3, &result.Trunk4, &result.Carrier4, &result.PlanCode4, &result.Plan4, &result.Duration4, &result.Trunk5, &result.Carrier5, &result.PlanCode5, &result.Plan5, &result.Duration5)
	if err != nil {
		return sharedData.RadiusOnestageValidateData{}, err
	}

	return result, nil
}
