package onevoisdata

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestExecuteRadiusAccounting(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	radiusAccountingData := RadiusAccountingDataModel{DB: db}
	input := RadiusAccounting{
		ConfID:       123,
		AccessNo:     "1234567890",
		Anino:        "9876543210",
		DestNo:       "5555555555",
		SubscriberNo: "4444444444",
		Pwd:          "password123",
		SessionID:    "session123",
		CategoryID:   "456",
		StartTime:    "2022-01-01 10:00:00",
		TalkingTime:  "00:30:00",
		CallDuration: 180,
		ReleaseCode:  "RELEASED",
		InTrunkID:    789,
		OutTrunkID:   321,
		ReasonID:     2,
		Prefix:       "prefix123",
		LanguageCode: "en",
	}

	rows := sqlmock.NewRows([]string{"STATUS", "CONF_ID", "REMAINMONEY"}).
		AddRow("SUCCESS", 123, 100.0)

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

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(input.ConfID,
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
			input.LanguageCode).
		WillReturnRows(rows)

	result, err := radiusAccountingData.ExecuteRadiusAccounting(input)
	assert.Nil(t, err)
	assert.Equal(t, "SUCCESS", result.Status)
	assert.Equal(t, 123, result.ConfID)
	assert.Equal(t, 100.0, result.RemainMoney)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestExecuteRadiusAccountingError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	radiusAccountingData := RadiusAccountingDataModel{DB: db}
	input := RadiusAccounting{
		ConfID:   123,
		AccessNo: "1234567890",
		// ... (other input fields)
	}

	mock.ExpectQuery(`.*`).
		WillReturnError(errors.New("database error"))

	_, err = radiusAccountingData.ExecuteRadiusAccounting(input)
	assert.Error(t, err)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}
