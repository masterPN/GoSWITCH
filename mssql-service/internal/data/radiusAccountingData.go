package data

import (
	"context"
	"database/sql"
	"time"
)

type RadiusAccountingData struct {
	Status      string `json:"STATUS"`
	ConfID      int    `json:"CONF_ID"`
	RemainMoney int    `json:"REMAINMONEY"`
}

type RadiusAccountingInput struct {
	ConfID       int    `json:"confID"`
	AccessNo     string `json:"accessNo"`
	Anino        string `json:"anino"`
	DestNo       string `json:"destNo"`
	SubscriberNo string `json:"subscriberNo"`
	Pwd          string `json:"pwd"`
	SessionID    string `json:"sessionID"`
	CategoryID   string `json:"categoryID"`
	StartTime    string `json:"startTime"`
	TalkingTime  string `json:"talkingTime"`
	CallDuration int    `json:"callDuration"`
	ReleaseCode  string `json:"releaseCode"`
	InTrunkID    int    `json:"inTrunkID"`
	OutTrunkID   int    `json:"outTrunkID"`
	ReasonID     int    `json:"reasonID"`
	Prefix       string `json:"prefix"`
	LanguageCode string `json:"languageCode"`
}

type RadiusAccountingDataModel struct {
	DB *sql.DB
}

func (r RadiusAccountingDataModel) ExecuteRadiusAccounting(input RadiusAccountingInput) (RadiusAccountingData, error) {
	query := `EXEC RADIUS_ACCOUNTING
				@CONF_ID = $1,
				@ACCESS_NO = $2,
				@ANINO = $3,
				@DEST_NO = $4,
				@SUBSCRIBER_NO = $5,
				@PWD = $6,
				@SESSION_ID = $7,
				@CATEGORY_ID = $8,
				@START_TIME = $9,
				@TALKING_TIME = $10,
				@CALL_DURATION = $11,
				@RELEASE_CODE = $12,
				@IN_TRUNK_ID = $13,
				@OUT_TRUNK_ID = $14,
				@REASON_ID = $15,
				@PREFIX = $16,
				@LANGUAGE_CODE = $17;`

	args := []interface{}{
		input.ConfID,
		input.AccessNo,
		input.Anino,
		input.DestNo,
		input.SubscriberNo,
		input.Pwd,
		input.SessionID,
		input.CategoryID,
		input.StartTime,
		input.TalkingTime,
		input.CallDuration,
		input.ReleaseCode,
		input.InTrunkID,
		input.OutTrunkID,
		input.ReasonID,
		input.Prefix,
		input.LanguageCode,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := r.DB.QueryRowContext(ctx, query, args...)

	var result RadiusAccountingData
	err := row.Scan(&result.Status, &result.ConfID, &result.RemainMoney)
	if err != nil {
		return RadiusAccountingData{}, err
	}

	return result, nil
}
