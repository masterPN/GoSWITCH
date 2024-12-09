package onevoisdata

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestExecuteRadiusOnestageValidate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := RadiusOnestageValidateDataModel{DB: db}

	t.Run("test successful execution with valid input", func(t *testing.T) {
		mock.ExpectQuery(`EXEC RADIUS_ONESTAGE_VALIDATE`).WithArgs("prefixdestinationNumber", "prefix", "callingNumber", "destinationNumber").WillReturnRows(sqlmock.NewRows([]string{"Lcode", "Status", "RouteType", "CallType", "AccountNum", "PrefixNo", "Dnis", "Pin", "FollowOnCall", "Trunk1", "Carrier1", "PlanCode1", "Plan1", "Duration1", "Trunk2", "Carrier2", "PlanCode2", "Plan2", "Duration2", "Trunk3", "Carrier3", "PlanCode3", "Plan3", "Duration3", "Trunk4", "Carrier4", "PlanCode4", "Plan4", "Duration4", "Trunk5", "Carrier5", "PlanCode5", "Plan5", "Duration5"}).AddRow("LcodeValue", 1, "RouteTypeValue", 1, "AccountNumValue", "PrefixNoValue", "DnisValue", "PinValue", "FollowOnCallValue", 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1))

		result, err := r.ExecuteRadiusOnestageValidate("prefix", "callingNumber", "destinationNumber")
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("test error handling for invalid input", func(t *testing.T) {
		_, err := r.ExecuteRadiusOnestageValidate("", "", "")
		assert.Error(t, err)
	})

	t.Run("test error handling for database query errors", func(t *testing.T) {
		mock.ExpectQuery(`EXEC RADIUS_ONESTAGE_VALIDATE`).WithArgs("prefixdestinationNumber", "prefix", "callingNumber", "destinationNumber").WillReturnError(sql.ErrConnDone)

		_, err := r.ExecuteRadiusOnestageValidate("prefix", "callingNumber", "destinationNumber")
		assert.Error(t, err)
	})

	t.Run("test error handling for Scan errors", func(t *testing.T) {
		mock.ExpectQuery(`EXEC RADIUS_ONESTAGE_VALIDATE`).WithArgs("prefixdestinationNumber", "prefix", "callingNumber", "destinationNumber").WillReturnRows(sqlmock.NewRows([]string{"Lcode", "Status", "RouteType", "CallType", "AccountNum", "PrefixNo", "Dnis", "Pin", "FollowOnCall", "Trunk1", "Carrier1", "PlanCode1", "Plan1", "Duration1", "Trunk2", "Carrier2", "PlanCode2", "Plan2", "Duration2", "Trunk3", "Carrier3", "PlanCode3", "Plan3", "Duration3", "Trunk4", "Carrier4", "PlanCode4", "Plan4", "Duration4", "Trunk5", "Carrier5", "PlanCode5", "Plan5", "Duration5"}).AddRow("Lcode", "Status", "RouteType", "CallType", "AccountNum", "PrefixNo", "Dnis", "Pin", "FollowOnCall", "Trunk1", "Carrier1", "PlanCode1", "Plan1", "Duration1", "Trunk2", "Carrier2", "PlanCode2", "Plan2", "Duration2", "Trunk3", "Carrier3", "PlanCode3", "Plan3", "Duration3", "Trunk4", "Carrier4", "PlanCode4", "Plan4", "Duration4", "Trunk5", "Carrier5", "PlanCode5", "Plan5", "Duration5"))
		_, err := r.ExecuteRadiusOnestageValidate("prefix", "callingNumber", "destinationNumber")
		assert.Error(t, err)
	})
}

func TestExecuteRadiusOnestageValidateContextTimeout(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := RadiusOnestageValidateDataModel{DB: db}

	mock.ExpectQuery(`EXEC RADIUS_ONESTAGE_VALIDATE`).WithArgs("prefixdestinationNumber", "prefix", "callingNumber", "destinationNumber").WillReturnError(context.DeadlineExceeded)

	_, err = r.ExecuteRadiusOnestageValidate("prefix", "callingNumber", "destinationNumber")
	assert.Error(t, err)
}
