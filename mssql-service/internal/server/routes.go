package server

import (
	"fmt"
	"mssql-service/internal/data/onevoisdata"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)
	r.POST("/radiusOnestageValidate", s.ExecuteRadiusOnestageValidateHandler)
	r.POST("/radiusAccounting", s.ExecuteRadiusAccountingHandler)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) ExecuteRadiusOnestageValidateHandler(c *gin.Context) {
	var input struct {
		Prefix            string `json:"prefix"`
		CallingNumber     string `json:"callingNumber"`
		DestinationNumber string `json:"destinationNumber"`
	}
	c.BindJSON(&input)

	result, err := s.onevoisModels.RadiusData.ExecuteRadiusOnestageValidate(input.Prefix, input.CallingNumber, input.DestinationNumber)
	if err != nil {
		c.Error(fmt.Errorf("ExecuteRadiusOnestageValidateHandler with %q - %q", input, err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s *Server) ExecuteRadiusAccountingHandler(c *gin.Context) {
	var input onevoisdata.RadiusAccountingInput
	c.BindJSON(&input)

	result, err := s.onevoisModels.RadiusAccountingData.ExecuteRadiusAccounting(input)
	if err != nil {
		c.Error(fmt.Errorf("ExecuteRadiusAccountingHandler with %q - %q", input, err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
