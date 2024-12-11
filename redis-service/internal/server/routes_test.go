package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"redis-service/internal/data"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	sharedConstants "github.com/masterPN/GoSWITCH-shared/constants"
	"github.com/stretchr/testify/assert"

	mock_data "redis-service/internal/data/mocks"
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

func TestSaveRadiusAccountingDataHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRadiusAccountingData := mock_data.NewMockRadiusAccountingDataInterface(ctrl)
	s := &Server{
		models: data.Models{RadiusAccountingData: mockRadiusAccountingData},
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
				"confID":   123,
				"accessNo": "1234567890",
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
			name: "failed data saving",
			input: gin.H{
				"confID":   123,
				"accessNo": "1234567890",
			},
			statusCode: http.StatusBadRequest,
			err:        errors.New(sharedConstants.ErrorPhrase),
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
			c.Request, _ = http.NewRequest("POST", "/saveRadiusAccountingData", bytes.NewBuffer(jsonBytes))

			var inputData data.RadiusAccountingData
			if test.name != sharedConstants.FailedJsonBindingPhrase {
				json.Unmarshal(jsonBytes, &inputData)
				if test.err != nil {
					mockRadiusAccountingData.EXPECT().Set(inputData).Return(test.err)
				} else {
					mockRadiusAccountingData.EXPECT().Set(inputData).Return(nil)
				}
			}

			s.SaveRadiusAccountingDataHandler(c)

			assert.Equal(t, test.statusCode, w.Code)
		})
	}
}

func TestPopRadiusAccountingDataHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRadiusAccountingData := mock_data.NewMockRadiusAccountingDataInterface(ctrl)

	tests := []struct {
		name       string
		anino      string
		statusCode int
		err        error
	}{
		{
			name:       "successful pop",
			anino:      "12345",
			statusCode: http.StatusOK,
			err:        nil,
		},
		{
			name:       "pop fails",
			anino:      "12345",
			statusCode: http.StatusBadRequest,
			err:        errors.New(sharedConstants.ErrorPhrase),
		},
		{
			name:       "anino is empty",
			anino:      "",
			statusCode: http.StatusBadRequest,
			err:        errors.New(sharedConstants.ErrorPhrase),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if test.err != nil {
				mockRadiusAccountingData.EXPECT().Pop(test.anino).Return(data.RadiusAccountingData{}, test.err)
			} else {
				mockRadiusAccountingData.EXPECT().Pop(test.anino).Return(data.RadiusAccountingData{ConfID: 123}, nil)
			}

			// Create a new server with the mock model
			s := &Server{
				models: data.Models{RadiusAccountingData: mockRadiusAccountingData},
			}

			// Set the anino parameter
			c.Params = gin.Params{{Key: "anino", Value: test.anino}}

			// Call the PopRadiusAccountingDataHandler function
			s.PopRadiusAccountingDataHandler(c)

			// Check the status code
			assert.Equal(t, test.statusCode, w.Code)

			// Check the response body
			if test.err != nil {
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, test.err.Error(), resp["error"])
			} else {
				var resp data.RadiusAccountingData
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, data.RadiusAccountingData{ConfID: 123}, resp)
			}
		})
	}
}

func TestSetInternalCodemappingDataHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInternalCodemappingData := mock_data.NewMockInternalCodemappingDataInterface(ctrl)
	s := &Server{
		models: data.Models{InternalCodemappingData: mockInternalCodemappingData},
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
				"InternalCode": "123",
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
			err:        errors.New(sharedConstants.ErrorPhrase),
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
			c.Request, _ = http.NewRequest("POST", "/internalCodemappingData", bytes.NewBuffer(jsonBytes))

			var inputData data.InternalCodemappingData
			json.Unmarshal(jsonBytes, &inputData)
			if inputData.InternalCode != 0 || test.name == sharedConstants.EmptyInputPhrase {
				mockInternalCodemappingData.EXPECT().Set(inputData).Return(test.err)
			}

			s.SetInternalCodemappingDataHandler(c)

			assert.Equal(t, test.statusCode, w.Code)
		})
	}
}

func TestGetInternalCodemappingDataHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInternalCodemappingData := mock_data.NewMockInternalCodemappingDataInterface(ctrl)
	s := &Server{
		models: data.Models{InternalCodemappingData: mockInternalCodemappingData},
	}

	tests := []struct {
		name                    string
		internalCodeString      string
		internalCode            int
		internalCodemappingData data.InternalCodemappingData
		err                     error
		statusCode              int
	}{
		{
			name:               "successful data retrieval",
			internalCodeString: "123",
			internalCode:       123,
			internalCodemappingData: data.InternalCodemappingData{
				ID:           123,
				InternalCode: 123,
				OperatorCode: 456,
			},
			err:        nil,
			statusCode: http.StatusOK,
		},
		{
			name:                    "invalid internal code",
			internalCodeString:      "abc",
			internalCode:            0,
			internalCodemappingData: data.InternalCodemappingData{},
			err:                     nil,
			statusCode:              http.StatusBadRequest,
		},
		{
			name:                    "failed to retrieve data",
			internalCodeString:      "123",
			internalCode:            123,
			internalCodemappingData: data.InternalCodemappingData{},
			err:                     errors.New(sharedConstants.ErrorPhrase),
			statusCode:              http.StatusBadRequest,
		},
		{
			name:                    sharedConstants.EmptyInputPhrase,
			internalCodeString:      "",
			internalCode:            0,
			internalCodemappingData: data.InternalCodemappingData{},
			err:                     nil,
			statusCode:              http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = gin.Params{{Key: "internalCode", Value: test.internalCodeString}}

			if test.internalCode != 0 {
				mockInternalCodemappingData.EXPECT().Get(test.internalCode).Return(test.internalCodemappingData, test.err)
			}

			s.GetInternalCodemappingDataHandler(c)

			assert.Equal(t, test.statusCode, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			if test.err != nil {
				assert.NoError(t, err)
				assert.Equal(t, test.err.Error(), resp["error"])
			} else {
				var internalCodemappingData data.InternalCodemappingData
				err := json.Unmarshal(w.Body.Bytes(), &internalCodemappingData)
				assert.NoError(t, err)
				assert.Equal(t, test.internalCodemappingData, internalCodemappingData)
			}
		})
	}
}

func TestDeleteInternalCodemappingDataHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockData := mock_data.NewMockInternalCodemappingDataInterface(ctrl)
	server := &Server{
		models: data.Models{InternalCodemappingData: mockData},
	}

	tests := []struct {
		name       string
		input      map[string]any
		statusCode int
		err        error
	}{
		{
			name: "successful deletion",
			input: gin.H{
				"InternalCode": 123,
			},
			statusCode: http.StatusOK,
			err:        nil,
		},
		{
			name: sharedConstants.FailedJsonBindingPhrase,
			input: gin.H{
				"InternalCode": "abc",
			},
			statusCode: http.StatusInternalServerError,
			err:        errors.New("json: cannot unmarshal string into Go struct field InternalCodemappingData.InternalCode of type int"),
		},
		{
			name: "deletion error",
			input: gin.H{
				"InternalCode": 123,
			},
			statusCode: http.StatusBadRequest,
			err:        errors.New("deletion error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			jsonBytes, _ := json.Marshal(test.input)
			c.Request, _ = http.NewRequest("DELETE", "/internalCodemapping", bytes.NewBuffer(jsonBytes))

			var input data.InternalCodemappingData
			json.Unmarshal(jsonBytes, &input)

			if input.InternalCode != 0 {
				mockData.EXPECT().Delete(input.InternalCode).Return(test.err)
			}

			server.DeleteInternalCodemappingDataHandler(c)

			assert.Equal(t, test.statusCode, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if test.err != nil {
				assert.Equal(t, test.err.Error(), response["error"])
			} else {
				assert.Equal(t, "Internal codemapping deleted successfully", response["message"])
			}
		})
	}
}
