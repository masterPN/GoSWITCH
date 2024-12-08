package wholesaledata

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetFirstImgCdrOperatorRoutingByNumber(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := ImgCdrOperatorRoutingDataModel{DB: db}

	t.Run("success", func(t *testing.T) {
		query := `SELECT TOP 1 * FROM IMGCDROPERATORROUTING IR WHERE $1 LIKE CONCAT(IR.ROUTINGPREFIX, '%') ORDER BY IR.ROUTINGPREFIX DESC;`
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("1234567890").WillReturnRows(sqlmock.NewRows([]string{"RoutingPrefix", "BaseOperator", "BaseClass1", "BaseClass2", "BaseClass3", "BaseClass4", "BaseClassCreationDate", "CurrentClass1", "CurrentClass2", "CurrentClass3", "CurrentClass4", "CurrentClassLastUpdateDate", "CurrentClassRemarks", "RoutingManualOverride", "OverrideCreationDate", "OverrideCreatedBy", "OverrideReason", "BaseClassLastUpdateDate"}).AddRow("1234567890", "Operator1", 1, 2, 3, 4, time.Now(), 1, 2, 3, 4, time.Now(), "Remarks", "Override", time.Now(), "CreatedBy", "Reason", time.Now()))

		result, err := r.GetFirstImgCdrOperatorRoutingByNumber("1234567890")
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("db query error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT TOP 1 * FROM IMGCDROPERATORROUTING IR WHERE $1 LIKE CONCAT\(IR\.ROUTINGPREFIX, '%'\) ORDER BY IR\.ROUTINGPREFIX DESC`).WithArgs("1234567890").WillReturnError(sql.ErrConnDone)

		_, err := r.GetFirstImgCdrOperatorRoutingByNumber("1234567890")
		assert.Error(t, err)
	})

	t.Run("scan error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT TOP 1 * FROM IMGCDROPERATORROUTING IR WHERE $1 LIKE CONCAT\(IR\.ROUTINGPREFIX, '%'\) ORDER BY IR\.ROUTINGPREFIX DESC`).WithArgs("1234567890").WillReturnRows(sqlmock.NewRows([]string{"RoutingPrefix", "BaseOperator", "BaseClass1", "BaseClass2", "BaseClass3", "BaseClass4", "BaseClassCreationDate", "CurrentClass1", "CurrentClass2", "CurrentClass3", "CurrentClass4", "CurrentClassLastUpdateDate", "CurrentClassRemarks", "RoutingManualOverride", "OverrideCreationDate", "OverrideCreatedBy", "OverrideReason", "BaseClassLastUpdateDate"}).AddRow("1234567890", "Operator1", 1, 2, 3, 4, time.Now(), 1, 2, 3, 4, time.Now(), "Remarks", "Override", time.Now(), "CreatedBy", "Reason", time.Now()))

		// simulate scan error
		mock.ExpectQuery(`SELECT TOP 1 * FROM IMGCDROPERATORROUTING IR WHERE $1 LIKE CONCAT\(IR\.ROUTINGPREFIX, '%'\) ORDER BY IR\.ROUTINGPREFIX DESC`).WithArgs("1234567890").WillReturnRows(sqlmock.NewRows([]string{"RoutingPrefix", "BaseOperator", "BaseClass1", "BaseClass2", "BaseClass3", "BaseClass4", "BaseClassCreationDate", "CurrentClass1", "CurrentClass2", "CurrentClass3", "CurrentClass4", "CurrentClassLastUpdateDate", "CurrentClassRemarks", "RoutingManualOverride", "OverrideCreationDate", "OverrideCreatedBy", "OverrideReason", "BaseClassLastUpdateDate"}).AddRow("1234567890", "Operator1", 1, 2, 3, 4, time.Now(), 1, 2, 3, 4, time.Now(), "Remarks", "Override", time.Now(), "CreatedBy", "Reason", time.Now()).RowError(1, sql.ErrNoRows))

		_, err := r.GetFirstImgCdrOperatorRoutingByNumber("1234567890")
		assert.Error(t, err)
	})

	t.Run("context timeout", func(t *testing.T) {
		_, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		mock.ExpectQuery(`SELECT TOP 1 * FROM IMGCDROPERATORROUTING IR WHERE $1 LIKE CONCAT\(IR\.ROUTINGPREFIX, '%'\) ORDER BY IR\.ROUTINGPREFIX DESC`).WithArgs("1234567890").WillDelayFor(2 * time.Millisecond)

		_, err := r.GetFirstImgCdrOperatorRoutingByNumber("1234567890")
		assert.Error(t, err)
	})
}
