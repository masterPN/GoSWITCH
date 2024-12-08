package onevoisdata

import (
	"context"
	"database/sql"
	"time"
)

type RadiusAccountingData struct {
	Status      string  `json:"STATUS"`
	ConfID      int     `json:"CONF_ID"`
	RemainMoney float64 `json:"REMAINMONEY"`
}

type RadiusAccounting struct {
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

func (r RadiusAccountingDataModel) ExecuteRadiusAccounting(input RadiusAccounting) (RadiusAccountingData, error) {
	query := `EXEC RADIUS_ACCOUNTING
                @CONF_ID = ?,
                @ACCESS_NO = ?,
                @ANINO = ?,
                @DEST_NO = ?,
                @SUBSCRIBER_NO = ?,
                @PWD = ?,
                @SESSION_ID = ?,
                @CATEGORY_ID = ?,
                @START_TIME = ?,
                @TALKING_TIME = ?,
                @CALL_DURATION = ?,
                @RELEASE_CODE = ?,
                @IN_TRUNK_ID = ?,
                @OUT_TRUNK_ID = ?,
                @REASON_ID = ?,
                @PREFIX = ?,
                @LANGUAGE_CODE = ?;`

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
