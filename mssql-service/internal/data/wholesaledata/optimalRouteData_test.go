package wholesaledata

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestOptimalRouteDataModelExecuteGetOptimalRoute(t *testing.T) {
	t.Run("successful execution", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta("EXEC GetOptimalRoute @pCallString = $1")).
			WithArgs("123").
			WillReturnRows(sqlmock.NewRows([]string{"Class1", "Class2", "Class3"}).
				AddRow(1, 2, 3))

		model := OptimalRouteDataModel{DB: db}
		result, err := model.ExecuteGetOptimalRoute("123")

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if result.Class1 != 1 || result.Class2 != 2 || result.Class3 != 3 {
			t.Errorf("wrong result: %v", result)
		}
	})

	t.Run("error executing query", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta("EXEC GetOptimalRoute @pCallString = $1")).
			WithArgs("123").
			WillReturnError(errors.New("some error"))

		model := OptimalRouteDataModel{DB: db}
		_, err = model.ExecuteGetOptimalRoute("123")

		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}
