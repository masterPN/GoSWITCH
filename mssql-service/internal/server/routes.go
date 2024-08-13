package server

import (
	"mssql-service/internal/data"
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

	result, err := s.models.RadiusData.ExecuteRadiusOnestageValidate(input.Prefix, input.CallingNumber, input.DestinationNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s *Server) ExecuteRadiusAccountingHandler(c *gin.Context) {
	var input data.RadiusAccountingInput
	c.BindJSON(&input)

	result, err := s.models.RadiusAccountingData.ExecuteRadiusAccounting(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
