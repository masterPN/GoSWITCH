package data

import (
	"errors"
	"reflect"
	"strconv"
	"testing"

	"github.com/go-redis/redismock/v8"
)

func TestInternalCodemappingDataModelSetRedisHSetError(t *testing.T) {
	internalCode := 123
	input := InternalCodemappingData{
		ID:           456,
		InternalCode: internalCode,
		OperatorCode: 789,
	}

	// Mock Redis client with a set error
	mockRedisClient, mock := redismock.NewClientMock()
	mock.ExpectHSet(internalCodemappingKeyPrefix+strconv.Itoa(internalCode), map[string]interface{}{
		"ID":           "456",
		"InternalCode": "123",
		"OperatorCode": "789",
	}).SetErr(errors.New("redis: set error"))

	// Create a new InternalCodemappingDataModel with the mocked Redis client
	internalCodemappingDataModel := InternalCodemappingDataModel{
		DB: mockRedisClient,
	}

	// Call the Set method with valid input
	err := internalCodemappingDataModel.Set(input)

	// Assert that the error is not nil and contains the expected message
	if err == nil {
		t.Error("Expected an error, but got nil")
	}
}

func TestInternalCodemappingDataModelSetSuccess(t *testing.T) {
	internalCode := 123
	input := InternalCodemappingData{
		ID:           456,
		InternalCode: internalCode,
		OperatorCode: 789,
	}

	// Mock Redis client with a successful set
	mockRedisClient, mock := redismock.NewClientMock()
	mock.ExpectHSet(internalCodemappingKeyPrefix+strconv.Itoa(internalCode), map[string]interface{}{
		"ID":           "456",
		"InternalCode": "123",
		"OperatorCode": "789",
	}).SetVal(1)

	// Create a new InternalCodemappingDataModel with the mocked Redis client
	internalCodemappingDataModel := InternalCodemappingDataModel{
		DB: mockRedisClient,
	}

	// Call the Set method with valid input
	err := internalCodemappingDataModel.Set(input)

	// Assert that the error is nil
	if err != nil {
		t.Errorf("Expected no error, but got %q", err.Error())
	}

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Not all expectations were met: %s", err.Error())
	}
}

func TestPopulateDataWithNonZeroIDInternalCodeOperatorCode(t *testing.T) {
	input := InternalCodemappingData{
		ID:           123,
		InternalCode: 456,
		OperatorCode: 789,
	}

	expectedData := map[string]interface{}{
		"ID":           "123",
		"InternalCode": "456",
		"OperatorCode": "789",
	}

	data := make(map[string]interface{})
	InternalCodemappingDataModel{}.populateData(data, &input)

	if !reflect.DeepEqual(data, expectedData) {
		t.Errorf("Expected data %v, but got %v", expectedData, data)
	}
}

func TestPopulateDataWithZeroIDInternalCodeOperatorCode(t *testing.T) {
	input := InternalCodemappingData{
		ID:           0,
		InternalCode: 0,
		OperatorCode: 0,
	}

	expectedData := make(map[string]interface{})

	data := make(map[string]interface{})
	InternalCodemappingDataModel{}.populateData(data, &input)

	if !reflect.DeepEqual(data, expectedData) {
		t.Errorf("Expected data %v, but got %v", expectedData, data)
	}
}

func TestGetInternalCodemappingDataRedisConnectionLost(t *testing.T) {
	internalCode := 123

	// Mock Redis client with a connection error
	mockRedisClient, mock := redismock.NewClientMock()
	mock.ExpectHGetAll(internalCodemappingKeyPrefix + strconv.Itoa(internalCode)).SetErr(errors.New("redis: client is closed"))

	// Create a new InternalCodemappingDataModel with the mocked Redis client
	internalCodemappingDataModel := InternalCodemappingDataModel{
		DB: mockRedisClient,
	}

	// Call the Get method with a valid internal code
	_, err := internalCodemappingDataModel.Get(internalCode)

	// Assert that the error is not nil and contains the expected message
	if err == nil {
		t.Error("Expected an error, but got nil")
	}
	expectedErrorMessage := "failed to get internal code mapping: redis: client is closed"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message %q, but got %q", expectedErrorMessage, err.Error())
	}
}

func TestGetInternalCodemappingDataRedisReturnsEmptyResult(t *testing.T) {
	internalCode := 123

	// Mock Redis client with empty result
	mockRedisClient, mock := redismock.NewClientMock()
	mock.ExpectHGetAll(internalCodemappingKeyPrefix + strconv.Itoa(internalCode)).SetVal(map[string]string{})

	// Create a new InternalCodemappingDataModel with the mocked Redis client
	internalCodemappingDataModel := InternalCodemappingDataModel{
		DB: mockRedisClient,
	}

	// Call the Get method with a valid internal code
	_, err := internalCodemappingDataModel.Get(internalCode)

	// Assert that the error is not nil and contains the expected message
	if err == nil {
		t.Error("Expected an error, but got nil")
	}
	expectedErrorMessage := "internal code mapping not found"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message %q, but got %q", expectedErrorMessage, err.Error())
	}
}

func TestGetInternalCodemappingDataRedisReturnsValidData(t *testing.T) {
	internalCode := 123

	// Mock Redis client with valid data
	mockRedisClient, mock := redismock.NewClientMock()
	internalCodeMappingData := map[string]string{
		"ID":           "123",
		"InternalCode": "123",
		"OperatorCode": "789",
	}
	mock.ExpectHGetAll(internalCodemappingKeyPrefix + strconv.Itoa(internalCode)).SetVal(internalCodeMappingData)

	// Create a new InternalCodemappingDataModel with the mocked Redis client
	internalCodemappingDataModel := InternalCodemappingDataModel{
		DB: mockRedisClient,
	}

	// Call the Get method with a valid internal code
	result, err := internalCodemappingDataModel.Get(internalCode)

	// Assert that the error is nil and the result matches the expected data
	if err != nil {
		t.Errorf("Expected no error, but got %q", err.Error())
	}
	expectedResult := InternalCodemappingData{
		ID:           123,
		InternalCode: 123,
		OperatorCode: 789,
	}
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected result %v, but got %v", expectedResult, result)
	}

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Not all expectations were met: %s", err.Error())
	}
}

func TestDeleteInternalCodemappingDataWithZeroInternalCode(t *testing.T) {
	internalCode := 0

	// Mock Redis client with a connection error
	mockRedisClient, _ := redismock.NewClientMock()

	// Create a new InternalCodemappingDataModel with the mocked Redis client
	internalCodemappingDataModel := InternalCodemappingDataModel{
		DB: mockRedisClient,
	}

	// Call the Delete method with a zero internal code
	err := internalCodemappingDataModel.Delete(internalCode)

	// Assert that the error is not nil and contains the expected message
	if err == nil {
		t.Error("Expected an error, but got nil")
	}
	expectedErrorMessage := "invalid internal code: 0"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message %q, but got %q", expectedErrorMessage, err.Error())
	}
}

func TestGetInternalCodemappingDataRedisError(t *testing.T) {
	internalCode := 123

	// Mock Redis client with an error
	mockRedisClient, mock := redismock.NewClientMock()
	mock.ExpectDel(internalCodemappingKeyPrefix + strconv.Itoa(internalCode)).SetErr(errors.New("redis: client is closed"))

	// Create a new InternalCodemappingDataModel with the mocked Redis client
	internalCodemappingDataModel := InternalCodemappingDataModel{
		DB: mockRedisClient,
	}

	// Call the Delete method with a valid internal code
	err := internalCodemappingDataModel.Delete(internalCode)

	// Assert that the error is not nil and contains the expected message
	if err == nil {
		t.Error("Expected an error, but got nil")
	}
	expectedErrorMessage := "failed to delete internal code mapping: redis: client is closed"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message %q, but got %q", expectedErrorMessage, err.Error())
	}
}

func TestDeleteInternalCodemappingDataSuccess(t *testing.T) {
	internalCode := 123

	// Mock Redis client with a successful deletion
	mockRedisClient, mock := redismock.NewClientMock()
	mock.ExpectDel(internalCodemappingKeyPrefix + strconv.Itoa(internalCode)).SetVal(1)

	// Create a new InternalCodemappingDataModel with the mocked Redis client
	internalCodemappingDataModel := InternalCodemappingDataModel{
		DB: mockRedisClient,
	}

	// Call the Delete method with a valid internal code
	err := internalCodemappingDataModel.Delete(internalCode)

	// Assert that the error is nil
	if err != nil {
		t.Errorf("Expected no error, but got %q", err.Error())
	}

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Not all expectations were met: %s", err.Error())
	}
}
