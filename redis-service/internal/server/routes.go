package server

import (
	"net/http"
	"redis-service/internal/data"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)
	r.POST("/saveRadiusAccountingData", s.SaveRadiusAccountingDataHandler)
	r.GET("/popRadiusAccountingData/:anino", s.PopRadiusAccountingDataHandler)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) SaveRadiusAccountingDataHandler(c *gin.Context) {
	var input data.RadiusAccountingData
	c.BindJSON(&input)

	err := s.models.RadiusAccountingData.Set(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *Server) PopRadiusAccountingDataHandler(c *gin.Context) {
	anino := c.Param("anino")

	radiusAccountingData, err := s.models.RadiusAccountingData.Pop(anino)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":                err.Error(),
			"radiusAccountingData": radiusAccountingData,
		})
		return
	}

	c.JSON(http.StatusOK, radiusAccountingData)
}
