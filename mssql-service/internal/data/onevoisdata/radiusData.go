package onevoisdata

import (
	"context"
	"database/sql"
	"time"
)

type RadiusData struct {
	Lcode        string `json:"LCODE"`
	Status       int    `json:"STATUS"`
	RouteType    string `json:"ROUTE_TYPE"`
	CallType     int    `json:"CALL_TYPE"`
	AccountNum   string `json:"ACCOUNT_NUMBER"`
	PrefixNo     string `json:"PREFIX_NO"`
	Dnis         string `json:"DNIS"`
	Pin          string `json:"PIN"`
	FollowOnCall string `json:"FOLLOWONCALL"`
	Trunk1       int    `json:"TRUNK1"`
	Carrier1     int    `json:"CARRIER1"`
	PlanCode1    int    `json:"PLAN_CODE1"`
	Plan1        int    `json:"PLAN1"`
	Duration1    int    `json:"DURATION1"`
	Trunk2       int    `json:"TRUNK2"`
	Carrier2     int    `json:"CARRIER2"`
	PlanCode2    int    `json:"PLAN_CODE2"`
	Plan2        int    `json:"PLAN2"`
	Duration2    int    `json:"DURATION2"`
	Trunk3       int    `json:"TRUNK3"`
	Carrier3     int    `json:"CARRIER3"`
	PlanCode3    int    `json:"PLAN_CODE3"`
	Plan3        int    `json:"PLAN3"`
	Duration3    int    `json:"DURATION3"`
	Trunk4       int    `json:"TRUNK4"`
	Carrier4     int    `json:"CARRIER4"`
	PlanCode4    int    `json:"PLAN_CODE4"`
	Plan4        int    `json:"PLAN4"`
	Duration4    int    `json:"DURATION4"`
	Trunk5       int    `json:"TRUNK5"`
	Carrier5     int    `json:"CARRIER5"`
	PlanCode5    int    `json:"PLAN_CODE5"`
	Plan5        int    `json:"PLAN5"`
	Duration5    int    `json:"DURATION5"`
}

type RadiusDataModel struct {
	DB *sql.DB
}

func (r RadiusDataModel) ExecuteRadiusOnestageValidate(prefix string, callingNumber string, destinationNumber string) (RadiusData, error) {
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

	var result RadiusData
	err := row.Scan(&result.Lcode, &result.Status, &result.RouteType, &result.CallType, &result.AccountNum, &result.PrefixNo, &result.Dnis, &result.Pin, &result.FollowOnCall, &result.Trunk1, &result.Carrier1, &result.PlanCode1, &result.Plan1, &result.Duration1, &result.Trunk2, &result.Carrier2, &result.PlanCode2, &result.Plan2, &result.Duration2, &result.Trunk3, &result.Carrier3, &result.PlanCode3, &result.Plan3, &result.Duration3, &result.Trunk4, &result.Carrier4, &result.PlanCode4, &result.Plan4, &result.Duration4, &result.Trunk5, &result.Carrier5, &result.PlanCode5, &result.Plan5, &result.Duration5)
	if err != nil {
		return RadiusData{}, err
	}

	return result, nil
}
