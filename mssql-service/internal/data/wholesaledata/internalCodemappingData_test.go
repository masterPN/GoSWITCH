package wholesaledata

import (
	"database/sql"
	"errors"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/masterPN/GoSWITCH-shared/constants"
	"github.com/stretchr/testify/assert"
)

const (
	selectInternalCodeMappingQuery       = "SELECT ID, InternalCode, OperatorCode FROM InternalCodeMapping"
	selectIdFromInternalCodeMappingQuery = "SELECT ID FROM InternalCodeMapping WHERE InternalCode = $1"
	deleteInternalCodeMappingQuery       = "DELETE FROM InternalCodeMapping WHERE InternalCode = $1"
)

func TestInternalCodemappingDataModelGetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	t.Run("Test Get successful retrieval of data", func(t *testing.T) {
		expectedRows := sqlmock.NewRows([]string{"ID", "InternalCode", "OperatorCode"}).
			AddRow(1, 123, 456).
			AddRow(2, 789, 101)
		mock.ExpectQuery(regexp.QuoteMeta(selectInternalCodeMappingQuery)).
			WillReturnRows(expectedRows)

		model := InternalCodemappingDataModel{DB: db}
		actualResults, err := model.GetAll()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expectedResults := []InternalCodemappingData{
			{ID: 1, InternalCode: 123, OperatorCode: 456},
			{ID: 2, InternalCode: 789, OperatorCode: 101},
		}
		if !reflect.DeepEqual(actualResults, expectedResults) {
			t.Errorf("Expected: %v, got: %v", expectedResults, actualResults)
		}
	})

	t.Run("Test Get error during query execution", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(selectInternalCodeMappingQuery)).
			WillReturnError(errors.New(constants.ErrorPhrase))
		model := InternalCodemappingDataModel{DB: db}

		_, err = model.GetAll()
		if err == nil {
			t.Error(constants.ExpectErrorPhrase)
		}
	})

	t.Run("Test Get error during row scanning", func(t *testing.T) {
		invalidRows := sqlmock.NewRows([]string{"ID", "InternalCode", "OperatorCode"}).
			AddRow(1, 123, 456).
			AddRow("invalid", 789, 101)
		mock.ExpectQuery(regexp.QuoteMeta(selectInternalCodeMappingQuery)).
			WillReturnRows(invalidRows)
		model := InternalCodemappingDataModel{DB: db}

		_, err = model.GetAll()
		if err == nil {
			t.Error(constants.ExpectErrorPhrase)
		}
	})

	t.Run("Test Get error retrieving rows", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(selectInternalCodeMappingQuery)).
			WillReturnError(sql.ErrNoRows)
		model := InternalCodemappingDataModel{DB: db}

		_, err = model.GetAll()
		if err == nil {
			t.Error(constants.ExpectErrorPhrase)
		}
	})
}

func TestInternalCodemappingDataModelSet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	t.Run("Test Set with valid internal code", func(t *testing.T) {
		input := InternalCodemappingData{
			InternalCode: 123,
			OperatorCode: 456,
		}

		mock.ExpectQuery(regexp.QuoteMeta(selectIdFromInternalCodeMappingQuery)).WithArgs(input.InternalCode).WillReturnRows(sqlmock.NewRows([]string{"ID"}))
		mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO InternalCodeMapping (InternalCode, OperatorCode) OUTPUT INSERTED.ID VALUES ($1, $2)")).WithArgs(input.InternalCode, input.OperatorCode).WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(1))
		model := InternalCodemappingDataModel{DB: db}
		result, err := model.Set(input)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.ID)
	})

	t.Run("Test Set with existing internal code", func(t *testing.T) {
		input := InternalCodemappingData{
			InternalCode: 123,
			OperatorCode: 789,
		}

		mock.ExpectQuery(regexp.QuoteMeta(selectIdFromInternalCodeMappingQuery)).WithArgs(input.InternalCode).WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(1))
		mock.ExpectQuery(regexp.QuoteMeta("UPDATE InternalCodeMapping SET InternalCode = $1, OperatorCode = $2 OUTPUT INSERTED.ID WHERE ID = $3")).WithArgs(input.InternalCode, input.OperatorCode, 1).WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(1))
		model := InternalCodemappingDataModel{DB: db}
		result, err := model.Set(input)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.ID)
	})

	t.Run("Test Set with error during query execution", func(t *testing.T) {
		input := InternalCodemappingData{
			InternalCode: 123,
			OperatorCode: 456,
		}

		mock.ExpectQuery(regexp.QuoteMeta(selectIdFromInternalCodeMappingQuery)).WithArgs(input.InternalCode).WillReturnError(errors.New(constants.ErrorPhrase))
		model := InternalCodemappingDataModel{DB: db}
		_, err := model.Set(input)

		assert.EqualError(t, err, constants.ErrorPhrase)
	})

	t.Run("Test Set with error during insert", func(t *testing.T) {
		input := InternalCodemappingData{
			InternalCode: 123,
			OperatorCode: 456,
		}

		mock.ExpectQuery(regexp.QuoteMeta(selectIdFromInternalCodeMappingQuery)).WithArgs(input.InternalCode).WillReturnError(sql.ErrNoRows)
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO InternalCodeMapping (InternalCode, OperatorCode) OUTPUT INSERTED.ID VALUES ($1, $2)")).WithArgs(input.InternalCode, input.OperatorCode).WillReturnError(sql.ErrConnDone)
		model := InternalCodemappingDataModel{DB: db}

		_, err = model.Set(input)

		if err == nil {
			t.Error(constants.ExpectErrorPhrase)
		}
	})

	t.Run("Test Set with error during update", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		input := InternalCodemappingData{
			InternalCode: 123,
			OperatorCode: 456,
		}

		mock.ExpectQuery(regexp.QuoteMeta(selectIdFromInternalCodeMappingQuery)).WithArgs(input.InternalCode).WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(1))
		mock.ExpectQuery(regexp.QuoteMeta("UPDATE InternalCodeMapping SET InternalCode = $1, OperatorCode = $2 OUTPUT INSERTED.ID WHERE ID = $3")).WithArgs(input.InternalCode, input.OperatorCode, 1).WillReturnError(errors.New("update error"))
		model := InternalCodemappingDataModel{DB: db}
		_, err = model.Set(input)

		if err == nil {
			t.Error(constants.ExpectErrorPhrase)
		}
		expectedErrorMessage := "update error"
		if err.Error() != expectedErrorMessage {
			t.Errorf("Expected error message %q, but got %q", expectedErrorMessage, err.Error())
		}
	})
}

func TestInternalCodemappingDataModelDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	model := InternalCodemappingDataModel{DB: db}

	t.Run("Test Delete with valid internal code", func(t *testing.T) {
		internalCode := 123
		mock.ExpectExec(regexp.QuoteMeta(deleteInternalCodeMappingQuery)).WithArgs(internalCode).WillReturnResult(sqlmock.NewResult(1, 1))
		err := model.Delete(internalCode)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
	})

	t.Run("Test Delete with invalid internal code", func(t *testing.T) {
		internalCode := 0
		mock.ExpectExec(regexp.QuoteMeta(deleteInternalCodeMappingQuery)).WithArgs(internalCode).WillReturnError(sql.ErrNoRows)
		err := model.Delete(internalCode)
		if err == nil {
			t.Errorf("Expected error, but got nil")
		}
	})

	t.Run("Test Delete with error during query execution", func(t *testing.T) {
		internalCode := 123
		mock.ExpectExec(regexp.QuoteMeta(deleteInternalCodeMappingQuery)).WithArgs(internalCode).WillReturnError(sql.ErrConnDone)
		err := model.Delete(internalCode)
		if err == nil {
			t.Errorf("Expected error, but got nil")
		}
	})
}
