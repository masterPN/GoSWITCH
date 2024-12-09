package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mock_data "mssql-service/internal/data/mocks"
	"mssql-service/internal/data/onevoisdata"
	"mssql-service/internal/data/wholesaledata"

	sharedConstants "github.com/masterPN/GoSWITCH-shared/constants"
	sharedData "github.com/masterPN/GoSWITCH-shared/data"
)

const (
	successfullyPhrase = "successful JSON binding and data execution"
)

func TestHelloWorldHandler(t *testing.T) {
	s := &Server{}
	r := gin.New()
	r.GET("/", s.HelloWorldHandler)
	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	// Serve the HTTP request
	r.ServeHTTP(rr, req)
	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	// Check the response body
	expected := "{\"message\":\"Hello World\"}"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestExecuteRadiusOnestageValidateHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRadiusOnestageValidateData := mock_data.NewMockRadiusOnestageValidateDataInterface(ctrl)
	s := &Server{
		onevoisModels: onevoisdata.Models{RadiusOnestageValidateData: mockRadiusOnestageValidateData},
	}

	tests := []struct {
		name       string
		input      map[string]any
		statusCode int
	}{
		{
			name: successfullyPhrase,
			input: gin.H{
				"prefix":            "prefix",
				"callingNumber":     "callingNumber",
				"destinationNumber": "destinationNumber",
			},
			statusCode: http.StatusOK,
		},
		{
			name: sharedConstants.FailedJsonBindingPhrase,
			input: gin.H{
				"prefix": 123,
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "error in ExecuteRadiusOnestageValidate function",
			input: gin.H{
				"prefix":            "prefix",
				"callingNumber":     "callingNumber",
				"destinationNumber": "destinationNumber",
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name:       sharedConstants.EmptyInputPhrase,
			input:      gin.H{},
			statusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			jsonBytes, _ := json.Marshal(test.input)
			c.Request, _ = http.NewRequest("POST", "/radiusOnestageValidate", bytes.NewBuffer(jsonBytes))

			var inputData sharedData.RadiusOnestageValidateData
			if test.name != sharedConstants.FailedJsonBindingPhrase {
				json.Unmarshal(jsonBytes, &inputData)
				if test.name == "error in ExecuteRadiusOnestageValidate function" {
					mockRadiusOnestageValidateData.EXPECT().ExecuteRadiusOnestageValidate("prefix", "callingNumber", "destinationNumber").Return(sharedData.RadiusOnestageValidateData{}, errors.New("error"))
				} else if test.name == sharedConstants.EmptyInputPhrase {
					mockRadiusOnestageValidateData.EXPECT().ExecuteRadiusOnestageValidate("", "", "").Return(sharedData.RadiusOnestageValidateData{}, nil)
				} else {
					mockRadiusOnestageValidateData.EXPECT().ExecuteRadiusOnestageValidate("prefix", "callingNumber", "destinationNumber").Return(sharedData.RadiusOnestageValidateData{}, nil)
				}
			}

			s.ExecuteRadiusOnestageValidateHandler(c)

			assert.Equal(t, test.statusCode, w.Code)
		})
	}
}

func TestExecuteRadiusAccountingHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRadiusAccountingData := mock_data.NewMockRadiusAccountingDataInterface(ctrl)
	s := &Server{
		onevoisModels: onevoisdata.Models{RadiusAccountingData: mockRadiusAccountingData},
	}

	tests := []struct {
		name       string
		input      map[string]any
		statusCode int
		err        error
	}{
		{
			name: successfullyPhrase,
			input: gin.H{
				"confID":       123,
				"accessNo":     "1234567890",
				"anino":        "9876543210",
				"destNo":       "5555555555",
				"subscriberNo": "4444444444",
				"pwd":          "password123",
				"sessionID":    "session123",
				"categoryID":   "456",
				"startTime":    "2022-01-01 10:00:00",
				"talkingTime":  "00:30:00",
				"callDuration": 180,
				"releaseCode":  "RELEASED",
				"inTrunkID":    789,
				"outTrunkID":   321,
				"reasonID":     2,
				"prefix":       "prefix123",
				"languageCode": "en",
			},
			statusCode: http.StatusOK,
		},
		{
			name: sharedConstants.FailedJsonBindingPhrase,
			input: gin.H{
				"confID": "123",
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "error in ExecuteRadiusAccounting function",
			input: gin.H{
				"confID":       123,
				"accessNo":     "1234567890",
				"anino":        "9876543210",
				"destNo":       "5555555555",
				"subscriberNo": "4444444444",
				"pwd":          "password123",
				"sessionID":    "session123",
				"categoryID":   "456",
				"startTime":    "2022-01-01 10:00:00",
				"talkingTime":  "00:30:00",
				"callDuration": 180,
				"releaseCode":  "RELEASED",
				"inTrunkID":    789,
				"outTrunkID":   321,
				"reasonID":     2,
				"prefix":       "prefix123",
				"languageCode": "en",
			},
			statusCode: http.StatusBadRequest,
			err:        errors.New(sharedConstants.ErrorPhrase),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			jsonBytes, _ := json.Marshal(test.input)
			c.Request, _ = http.NewRequest("POST", "/radiusAccounting", bytes.NewBuffer(jsonBytes))

			var inputData onevoisdata.RadiusAccounting
			if test.name != sharedConstants.FailedJsonBindingPhrase {
				json.Unmarshal(jsonBytes, &inputData)
				if test.err != nil {
					mockRadiusAccountingData.EXPECT().ExecuteRadiusAccounting(inputData).Return(onevoisdata.RadiusAccountingData{}, test.err)
				} else {
					mockRadiusAccountingData.EXPECT().ExecuteRadiusAccounting(inputData).Return(onevoisdata.RadiusAccountingData{}, nil)
				}
			}

			s.ExecuteRadiusAccountingHandler(c)

			assert.Equal(t, test.statusCode, w.Code)
		})
	}
}

func TestGetOperatorRoutingHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOperatorRoutingData := mock_data.NewMockImgCdrOperatorRoutingDataInterface(ctrl)
	s := &Server{
		wholesaleModels: wholesaledata.Models{ImgCdrOperatorRoutingData: mockOperatorRoutingData},
	}

	tests := []struct {
		name       string
		number     string
		statusCode int
		err        error
	}{
		{
			name:       successfullyPhrase,
			number:     "1234567890",
			statusCode: http.StatusOK,
		},
		{
			name:       "error in GetOperatorRouting function",
			number:     "1234567890",
			statusCode: http.StatusBadRequest,
			err:        errors.New(sharedConstants.ErrorPhrase),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req, _ := http.NewRequest("GET", "/operatorRouting", nil)

			q := req.URL.Query()
			q.Add("number", test.number)
			req.URL.RawQuery = q.Encode()
			c.Request = req

			if test.name != sharedConstants.FailedJsonBindingPhrase {
				if test.err != nil {
					mockOperatorRoutingData.EXPECT().GetFirstImgCdrOperatorRoutingByNumber(test.number).Return(wholesaledata.ImgCdrOperatorRoutingData{}, test.err)
				} else {
					mockOperatorRoutingData.EXPECT().GetFirstImgCdrOperatorRoutingByNumber(test.number).Return(wholesaledata.ImgCdrOperatorRoutingData{}, nil)
				}
			}

			s.GetOperatorRoutingHandler(c)

			assert.Equal(t, test.statusCode, w.Code)
		})
	}
}

func TestExecuteGetOptimalRouteHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOptimalRouteData := mock_data.NewMockOptimalRouteDataInterface(ctrl)
	s := &Server{
		wholesaleModels: wholesaledata.Models{OptimalRouteData: mockOptimalRouteData},
	}

	tests := []struct {
		name           string
		pCallString    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "valid pCallString",
			pCallString:    "123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid pCallString",
			pCallString:    "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "error from OptimalRouteData.ExecuteGetOptimalRoute",
			pCallString:    "123",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req, _ := http.NewRequest("GET", "/optimalRoute", nil)

			q := req.URL.Query()
			q.Add("pCallString", test.pCallString)
			req.URL.RawQuery = q.Encode()
			c.Request = req

			if test.name != "valid pCallString" {
				mockOptimalRouteData.EXPECT().ExecuteGetOptimalRoute(test.pCallString).Return(wholesaledata.OptimalRouteData{}, errors.New("some error"))
			} else {
				mockOptimalRouteData.EXPECT().ExecuteGetOptimalRoute(test.pCallString).Return(wholesaledata.OptimalRouteData{}, nil)
			}

			s.ExecuteGetOptimalRouteHandler(c)

			assert.Equal(t, test.expectedStatus, w.Code)

			if test.expectedBody != "" {
				var body map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, test.expectedBody, body["error"])
			}
		})
	}
}

func TestFetchAllInternalCodemappingsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInternalCodemappingData := mock_data.NewMockInternalCodemappingDataInterface(ctrl)
	s := &Server{
		wholesaleModels: wholesaledata.Models{InternalCodemappingData: mockInternalCodemappingData},
	}

	tests := []struct {
		name       string
		err        error
		statusCode int
	}{
		{
			name:       "successful retrieval",
			err:        nil,
			statusCode: http.StatusOK,
		},
		{
			name:       "retrieval error",
			err:        errors.New("test error"),
			statusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			mockInternalCodemappingData.EXPECT().GetAll().Return(nil, test.err)

			s.FetchAllInternalCodemappingsHandler(c)

			if w.Code != test.statusCode {
				t.Errorf("expected status code %d, got %d", test.statusCode, w.Code)
			}

			if test.err != nil {
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				if err != nil {
					t.Errorf("failed to unmarshal response: %v", err)
				}
				if resp["error"] != test.err.Error() {
					t.Errorf("expected error message %q, got %q", test.err.Error(), resp["error"])
				}
			}
		})
	}
}

func TestSetInternalCodemappingHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInternalCodemappingData := mock_data.NewMockInternalCodemappingDataInterface(ctrl)
	s := &Server{
		wholesaleModels: wholesaledata.Models{InternalCodemappingData: mockInternalCodemappingData},
	}

	tests := []struct {
		name       string
		input      map[string]any
		statusCode int
		err        error
	}{
		{
			name: "successful JSON binding and data saving",
			input: gin.H{
				"InternalCode": 123,
				"OperatorCode": 456,
			},
			statusCode: http.StatusOK,
		},
		{
			name: sharedConstants.FailedJsonBindingPhrase,
			input: gin.H{
				"InternalCode": "abc",
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "failed data saving",
			input: gin.H{
				"InternalCode": 123,
				"OperatorCode": 456,
			},
			statusCode: http.StatusBadRequest,
			err:        errors.New("error saving data"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			jsonBytes, _ := json.Marshal(test.input)
			c.Request, _ = http.NewRequest("POST", "/internalCodemapping", bytes.NewBuffer(jsonBytes))

			if test.name != sharedConstants.FailedJsonBindingPhrase {
				if test.err != nil {
					mockInternalCodemappingData.EXPECT().Set(gomock.Any()).Return(wholesaledata.InternalCodemappingData{}, test.err)
				} else {
					mockInternalCodemappingData.EXPECT().Set(gomock.Any()).Return(wholesaledata.InternalCodemappingData{}, nil)
				}
			}

			s.SetInternalCodemappingHandler(c)

			assert.Equal(t, test.statusCode, w.Code)
		})
	}
}

func TestDeleteInternalCodemappingHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInternalCodemappingData := mock_data.NewMockInternalCodemappingDataInterface(ctrl)
	s := &Server{
		wholesaleModels: wholesaledata.Models{
			InternalCodemappingData: mockInternalCodemappingData,
		},
	}

	tests := []struct {
		name               string
		internalCodeString string
		err                error
		statusCode         int
	}{
		{
			name:               "successful deletion",
			internalCodeString: "123",
			err:                nil,
			statusCode:         http.StatusOK,
		},
		{
			name:               "error converting internal code string to integer",
			internalCodeString: "abc",
			err:                errors.New("strconv.Atoi: parsing \"abc\": invalid syntax"),
			statusCode:         http.StatusInternalServerError,
		},
		{
			name:               "error deleting internal codemapping data",
			internalCodeString: "123",
			err:                errors.New("deletion error"),
			statusCode:         http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req, _ := http.NewRequest("DELETE", "/internalCodemapping", nil)

			q := req.URL.Query()
			q.Add("internalCode", test.internalCodeString)
			req.URL.RawQuery = q.Encode()
			c.Request = req

			if test.name != "error converting internal code string to integer" {
				mockInternalCodemappingData.EXPECT().Delete(gomock.Any()).Return(test.err)
			}

			s.DeleteInternalCodemappingHandler(c)

			assert.Equal(t, test.statusCode, w.Code)

			if test.err != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, test.err.Error(), response["error"])
			}
		})
	}
}
