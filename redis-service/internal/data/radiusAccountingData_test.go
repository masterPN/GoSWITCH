package data

import (
	"errors"
	"reflect"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	sharedConstants "github.com/masterPN/GoSWITCH-shared/constants"
)

func TestSetFunctionHandlesNilInputGracefully(t *testing.T) {
	// Create a mock Redis client
	mockRedisClient := &redis.Client{}

	// Create a RadiusAccountingDataModel with the mock Redis client
	radiusAccountingDataModel := RadiusAccountingDataModel{
		DB: mockRedisClient,
	}

	// Call the Set function with a nil input
	err := radiusAccountingDataModel.Set(RadiusAccountingData{})

	// Assert that the error is nil
	if err != nil {
		t.Errorf(sharedConstants.ExpectedErrorPhrase, err)
	}
}

func TestSetFunctionWithNonEmptyFields(t *testing.T) {
	// Create a mock Redis client
	mockRedisClient, mock := redismock.NewClientMock()

	// Create a RadiusAccountingDataModel with the mock Redis client
	radiusAccountingDataModel := RadiusAccountingDataModel{
		DB: mockRedisClient,
	}

	// Define the input RadiusAccountingData with non-empty fields
	input := RadiusAccountingData{
		ConfID:       123,
		AccessNo:     "1234567890",
		StartTime:    sharedConstants.MockTimeString,
		CallDuration: 60,
		SessionID:    "01",
	}

	// Mock the HSet method to expect the input data
	mock.ExpectHSet(
		input.SessionID,
		map[string]interface{}{
			"ConfID":       int64(input.ConfID),
			"AccessNo":     input.AccessNo,
			"StartTime":    input.StartTime,
			"CallDuration": int64(input.CallDuration),
			"SessionID":    input.SessionID,
		},
	).SetVal(0)

	// Call the Set function with the input data
	err := radiusAccountingDataModel.Set(input)

	// Assert that the error is nil
	if err != nil {
		t.Errorf(sharedConstants.ExpectedErrorPhrase, err)
	}
}

func TestSetFunctionErrorDuringHSet(t *testing.T) {
	// Create a mock Redis client
	mockRedisClient, mock := redismock.NewClientMock()

	// Create a RadiusAccountingDataModel with the mock Redis client
	radiusAccountingDataModel := RadiusAccountingDataModel{
		DB: mockRedisClient,
	}

	// Define the input RadiusAccountingData with non-empty fields
	input := RadiusAccountingData{
		ConfID:       123,
		AccessNo:     "1234567890",
		StartTime:    sharedConstants.MockTimeString,
		CallDuration: 60,
		SessionID:    "01",
	}

	// Mock the HSet method to return an error
	mock.ExpectHSet(
		input.SessionID,
		map[string]interface{}{
			"ConfID":       int64(input.ConfID),
			"AccessNo":     input.AccessNo,
			"StartTime":    input.StartTime,
			"CallDuration": int64(input.CallDuration),
			"SessionID":    input.SessionID,
		},
	).SetErr(errors.New(sharedConstants.RedisClosedPhrase))

	// Call the Set function with the input data
	err := radiusAccountingDataModel.Set(input)

	// Assert that the error is not nil
	if err == nil {
		t.Error(sharedConstants.ExpectErrorPhrase)
	}

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(sharedConstants.MissExpectationPhrase, err)
	}
}

func TestPopulateDataHandlesZeroValuesForIntFields(t *testing.T) {
	r := RadiusAccountingDataModel{}

	input := RadiusAccountingData{
		ConfID:       0,
		CallDuration: 0,
		InTrunkID:    0,
		OutTrunkID:   0,
		ReasonID:     0,
	}

	data := make(map[string]interface{})
	r.populateData(data, input)

	expectedData := map[string]interface{}{}

	if !reflect.DeepEqual(data, expectedData) {
		t.Errorf("Expected data %v, but got %v", expectedData, data)
	}
}

func TestPopulateDataHandlesEmptyStrings(t *testing.T) {
	radiusAccountingData := RadiusAccountingData{
		ConfID:       123,
		AccessNo:     "",
		Anino:        "1234567890",
		DestNo:       "",
		SubscriberNo: "",
		Pwd:          "",
		SessionID:    "01",
		CategoryID:   "",
		StartTime:    sharedConstants.MockTimeString,
		TalkingTime:  "",
		CallDuration: 60,
		ReleaseCode:  "",
		InTrunkID:    1,
		OutTrunkID:   2,
		ReasonID:     3,
		Prefix:       "",
		LanguageCode: "",
	}

	expectedData := map[string]interface{}{
		"ConfID":       int64(123),
		"Anino":        "1234567890",
		"StartTime":    sharedConstants.MockTimeString,
		"CallDuration": int64(60),
		"SessionID":    "01",
		"InTrunkID":    int64(1),
		"OutTrunkID":   int64(2),
		"ReasonID":     int64(3),
	}

	data := make(map[string]interface{})
	radiusAccountingDataModel := RadiusAccountingDataModel{}
	radiusAccountingDataModel.populateData(data, radiusAccountingData)

	if !reflect.DeepEqual(data, expectedData) {
		t.Errorf("Expected data %v, but got %v", expectedData, data)
	}
}

func TestPopFunctionReturnsErrorWhenHGetAllReturnsError(t *testing.T) {
	// Create a mock Redis client
	mockRedisClient, mock := redismock.NewClientMock()

	// Set up a mock for the HGetAll method
	mock.ExpectHGetAll(
		"session123",
	).SetErr(
		errors.New(sharedConstants.RedisClosedPhrase),
	)

	// Create a RadiusAccountingDataModel with the mock Redis client
	radiusAccountingDataModel := RadiusAccountingDataModel{
		DB: mockRedisClient,
	}

	// Call the Pop function with a sessionID
	_, err := radiusAccountingDataModel.Pop("session123")

	// Assert that the error is not nil
	if err == nil {
		t.Error(sharedConstants.ExpectErrorPhrase)
	}
}

func TestPopFunctionReturnsErrorWhenDelReturnsError(t *testing.T) {
	// Create a mock Redis client
	mockRedisClient, mock := redismock.NewClientMock()

	// Set up a mock for the HGetAll method
	mock.ExpectHGetAll(
		"session123",
	).SetVal(
		map[string]string{
			"confID": "123",
			// Add other fields as needed
		},
	)

	// Set up a mock for the Del method to return an error
	mock.ExpectDel(
		"session123",
	).SetErr(
		errors.New(sharedConstants.RedisClosedPhrase),
	)

	// Create a RadiusAccountingDataModel with the mock Redis client
	radiusAccountingDataModel := RadiusAccountingDataModel{
		DB: mockRedisClient,
	}

	// Call the Pop function with a sessionID
	_, err := radiusAccountingDataModel.Pop("session123")

	// Assert that the error is not nil
	if err == nil {
		t.Error(sharedConstants.ExpectErrorPhrase)
	}

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(sharedConstants.MissExpectationPhrase, err)
	}
}

func TestPopFunctionReturnsCorrectRadiusAccountingData(t *testing.T) {
	// Create a mock Redis client
	mockRedisClient, mock := redismock.NewClientMock()

	// Set up a mock for the HGetAll method
	mock.ExpectHGetAll(
		"session123",
	).SetVal(
		map[string]string{
			"confID":       "123",
			"accessNo":     "1234567890",
			"anino":        "9876543210",
			"destNo":       "5555555555",
			"subscriberNo": "3333333333",
			"pwd":          "password",
			"sessionID":    "session123",
			"categoryID":   "category1",
			"startTime":    sharedConstants.MockTimeString,
			"talkingTime":  "00:01:00",
			"callDuration": "60",
			"releaseCode":  "release1",
			"inTrunkID":    "1",
			"outTrunkID":   "2",
			"reasonID":     "3",
			"prefix":       "prefix1",
			"languageCode": "en",
		},
	)

	// Set up a mock for the Del method
	mock.ExpectDel(
		"session123",
	).SetVal(
		1,
	)

	// Create a RadiusAccountingDataModel with the mock Redis client
	radiusAccountingDataModel := RadiusAccountingDataModel{
		DB: mockRedisClient,
	}

	// Call the Pop function with a sessionID
	radiusAccountingData, err := radiusAccountingDataModel.Pop("session123")

	// Assert that the error is nil
	if err != nil {
		t.Errorf(sharedConstants.ExpectedErrorPhrase, err)
	}

	// Assert that the returned RadiusAccountingData has the correct values
	expectedRadiusAccountingData := RadiusAccountingData{
		ConfID:       123,
		AccessNo:     "1234567890",
		Anino:        "9876543210",
		DestNo:       "5555555555",
		SubscriberNo: "3333333333",
		Pwd:          "password",
		SessionID:    "session123",
		CategoryID:   "category1",
		StartTime:    sharedConstants.MockTimeString,
		TalkingTime:  "00:01:00",
		CallDuration: 60,
		ReleaseCode:  "release1",
		InTrunkID:    1,
		OutTrunkID:   2,
		ReasonID:     3,
		Prefix:       "prefix1",
		LanguageCode: "en",
	}

	if !reflect.DeepEqual(radiusAccountingData, expectedRadiusAccountingData) {
		t.Errorf("Expected RadiusAccountingData %v, but got %v", expectedRadiusAccountingData, radiusAccountingData)
	}

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(sharedConstants.MissExpectationPhrase, err)
	}
}

func TestPopFunctionWithNotExistKey(t *testing.T) {
	// Create a mock Redis client
	mockRedisClient, mock := redismock.NewClientMock()

	// Set up a mock for the HGetAll method
	mock.ExpectHGetAll(
		"session123",
	).SetVal(
		map[string]string{
			"accessNo":     "1234567890",
			"startTime":    sharedConstants.MockTimeString,
			"callDuration": "60",
		},
	)

	// Set up a mock for the Del method
	mock.ExpectDel(
		"session123",
	).SetVal(
		0,
	)

	// Create a RadiusAccountingDataModel with the mock Redis client
	radiusAccountingDataModel := RadiusAccountingDataModel{
		DB: mockRedisClient,
	}

	// Call the Pop function with a sessionID
	radiusAccountingData, err := radiusAccountingDataModel.Pop("session123")

	// Assert that the error is nil
	if err == nil {
		t.Errorf("Expected error")
	}

	// Assert that the returned RadiusAccountingData has zero values for int fields
	expectedRadiusAccountingData := RadiusAccountingData{
		AccessNo:     "1234567890",
		StartTime:    sharedConstants.MockTimeString,
		CallDuration: 60,
		// Other fields have default values
	}

	if !reflect.DeepEqual(radiusAccountingData, expectedRadiusAccountingData) {
		t.Errorf("Expected RadiusAccountingData %v, but got %v", expectedRadiusAccountingData, radiusAccountingData)
	}

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(sharedConstants.MissExpectationPhrase, err)
	}
}
