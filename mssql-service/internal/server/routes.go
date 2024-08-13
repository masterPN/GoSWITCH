package server

import (
	"net/http"
	"time"

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
	var input struct {
		ConfID       int       `json:"confID"`
		AccessNo     string    `json:"accessNo"`
		Anino        string    `json:"anino"`
		DestNo       string    `json:"destNo"`
		SubscriberNo string    `json:"subscriberNo"`
		Pwd          string    `json:"pwd"`
		SessionID    string    `json:"sessionID"`
		CategoryID   string    `json:"categoryID"`
		StartTime    time.Time `json:"startTime"`
		TalkingTime  time.Time `json:"talkingTime"`
		CallDuration int       `json:"callDuration"`
		ReleaseCode  string    `json:"releaseCode"`
		InTrunkID    int       `json:"inTrunkID"`
		OutTrunkID   int       `json:"outTrunkID"`
		ReasonID     int       `json:"reasonID"`
		Prefix       string    `json:"prefix"`
		LanguageCode string    `json:"languageCode"`
	}
	c.BindJSON(&input)
}
