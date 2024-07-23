package data

import (
	"context"
	"database/sql"
	"time"
)

type CallLog struct {
	CallLogID        int       `json:"CALLLOG_ID"`
	SubscriberNo     string    `json:"SUBSCRIBER_NO"`
	StartTime        time.Time `json:"START_TIME"`
	TalkingTime      time.Time `json:"TALKING_TIME"`
	LegaDuration     int       `json:"LEGA_DURATION"`
	CallDuration     int       `json:"CALL_DURATION"`
	SysName          string    `json:"SYS_NAME"`
	ServerName       string    `json:"SERVER_NAME"`
	LanguageCode     string    `json:"LANGUAGE_CODE"`
	InLineNo         int       `json:"IN_LINE_NO"`
	OutLineNo        int       `json:"OUT_LINE_NO"`
	ReleaseCode      string    `json:"RELEASE_CODE"`
	ReasonID         int       `json:"REASON_ID"`
	AniNo            string    `json:"ANI_NO"`
	DestNo           string    `json:"DEST_NO"`
	AccessNo         string    `json:"ACCESS_NO"`
	CategoryID       string    `json:"CATEGORY_ID"`
	ServiceType      int       `json:"SERVICE_TYPE"`
	AuthenticateType int       `json:"AUTHENTICATE_TYPE"`
	SessionID        string    `json:"SESSION_ID"`
	InTrunkID        int       `json:"IN_TRUNK_ID"`
	OutTrunkID       int       `json:"OUT_TRUNK_ID"`
	InfoDigit        string    `json:"INFO_DIGIT"`
	InfoSurcharge    *float64  `json:"INFO_SURCHARGE"`
	ChargedADuration int       `json:"CHARGED_ADURATION"`
	ChargedDuration  int       `json:"CHARGED_DURATION"`
	LegaRate         string    `json:"LEGA_RATE"`
	LegbRate         string    `json:"LEGB_RATE"`
	CallACharge      string    `json:"CALL_ACHARGE"`
	CallBCharge      string    `json:"CALL_BCHARGE"`
	CallCharge       string    `json:"CALL_CHARGE"`
	RPlanID          int       `json:"RPLAN_ID"`
	AniTransCode     string    `json:"ANI_TRANS_CODE"`
	DnisTransCode    string    `json:"DNIS_TRANS_CODE"`
	ACountryCode     string    `json:"A_COUNTRY_CODE"`
	AAreaCode        string    `json:"A_AREA_CODE"`
	AAreaName        string    `json:"A_AREA_NAME"`
	BCountryCode     string    `json:"B_COUNTRY_CODE"`
	BAreaCode        string    `json:"B_AREA_CODE"`
	BAreaName        string    `json:"B_AREA_NAME"`
	ValueLeft        float64   `json:"VALUE_LEFT"`
	BatchNo          int       `json:"BATCH_NO"`
	ResellerID       string    `json:"RESELLER_ID"`
	LegaPrefix       string    `json:"LEGA_PREFIX"`
	RoutePrefix      string    `json:"ROUTE_PREFIX"`
	SPID             string    `json:"SP_ID"`
	LegaCostRate     string    `json:"LEGA_COST_RATE"`
	LegbCostRate     string    `json:"LEGB_COST_RATE"`
	LegaCostCharge   string    `json:"LEGA_COST_CHARGE"`
	LegbCostCharge   string    `json:"LEGB_COST_CHARGE"`
	LegaFRate        string    `json:"LEGA_F_RATE"`
	LegbFRate        string    `json:"LEGB_F_RATE"`
	LegaFDuration    int       `json:"LEGA_F_DURATION"`
	LegbFDuration    int       `json:"LEGB_F_DURATION"`
	LegaFCharge      string    `json:"LEGA_F_CHARGE"`
	LegbFCharge      string    `json:"LEGB_F_CHARGE"`
	LegaSRate        string    `json:"LEGA_S_RATE"`
	LegbSRate        string    `json:"LEGB_S_RATE"`
	LegaSDuration    int       `json:"LEGA_S_DURATION"`
	LegbSDuration    int       `json:"LEGB_S_DURATION"`
	LegaSCharge      string    `json:"LEGA_S_CHARGE"`
	LegbSCharge      string    `json:"LEGB_S_CHARGE"`
	LegaTRate        string    `json:"LEGA_T_RATE"`
	LegbTRate        string    `json:"LEGB_T_RATE"`
	LegaTDuration    int       `json:"LEGA_T_DURATION"`
	LegbTDuration    int       `json:"LEGB_T_DURATION"`
	LegaTCharge      string    `json:"LEGA_T_CHARGE"`
	LegbTCharge      string    `json:"LEGB_T_CHARGE"`
}

type CallLogModel struct {
	DB *sql.DB
}

func (c CallLogModel) Insert(callLog *CallLog) error {
	query := `
		INSERT INTO CALLLOG_1801
		(SUBSCRIBER_NO, START_TIME, TALKING_TIME, LEGA_DURATION, CALL_DURATION, SYS_NAME, SERVER_NAME, LANGUAGE_CODE, IN_LINE_NO, OUT_LINE_NO, RELEASE_CODE, REASON_ID, ANI_NO, DEST_NO, ACCESS_NO, CATEGORY_ID, SERVICE_TYPE, AUTHENTICATE_TYPE, SESSION_ID, IN_TRUNK_ID, OUT_TRUNK_ID, INFO_DIGIT, INFO_SURCHARGE, CHARGED_ADURATION, CHARGED_DURATION, LEGA_RATE, LEGB_RATE, CALL_ACHARGE, CALL_BCHARGE, CALL_CHARGE, RPLAN_ID, ANI_TRANS_CODE, DNIS_TRANS_CODE, A_COUNTRY_CODE, A_AREA_CODE, A_AREA_NAME, B_COUNTRY_CODE, B_AREA_CODE, B_AREA_NAME, VALUE_LEFT, BATCH_NO, RESELLER_ID, LEGA_PREFIX, ROUTE_PREFIX, SP_ID, LEGA_COST_RATE, LEGB_COST_RATE, LEGA_COST_CHARGE, LEGB_COST_CHARGE, LEGA_F_RATE, LEGB_F_RATE, LEGA_F_DURATION, LEGB_F_DURATION, LEGA_F_CHARGE, LEGB_F_CHARGE, LEGA_S_RATE, LEGB_S_RATE, LEGA_S_DURATION, LEGB_S_DURATION, LEGA_S_CHARGE, LEGB_S_CHARGE, LEGA_T_RATE, LEGB_T_RATE, LEGA_T_DURATION, LEGB_T_DURATION, LEGA_T_CHARGE, LEGB_T_CHARGE)
		VALUES ('', $1, $2, 0, 0, '', '', '', 0, 0, '00', 999, '', '', '', '', 0, 1, '', 0, 0, '', 0, 0, 0, 0, 0, 0, 0, 0, 0, '', '', '', '', '', '', '', '', 0, 0, '', '', '', '', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)`

	args := []interface{}{callLog.StartTime, callLog.TalkingTime}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.DB.ExecContext(ctx, query, args...)

	return err
}
