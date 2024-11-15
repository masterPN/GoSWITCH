package server

import (
	"fmt"
	"net/http"
	"redis-service/internal/data"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)
	r.POST("/saveRadiusAccountingData", s.SaveRadiusAccountingDataHandler)
	r.GET("/popRadiusAccountingData/:anino", s.PopRadiusAccountingDataHandler)
	r.POST("/setInternalCodemappingData", s.SetInternalCodemappingDataHandler)

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
		c.Error(fmt.Errorf("SaveRadiusAccountingDataHandler with %q - %q", input, err.Error()))
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
		c.Error(fmt.Errorf("PopRadiusAccountingDataHandler with %q - %q", anino, err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":                err.Error(),
			"radiusAccountingData": radiusAccountingData,
		})
		return
	}

	c.JSON(http.StatusOK, radiusAccountingData)
}

func (s *Server) SetInternalCodemappingDataHandler(c *gin.Context) {
	var input data.InternalCodemappingData
	c.BindJSON(&input)

	err := s.models.InternalCodemappingData.Set(input)
	if err != nil {
		c.Error(fmt.Errorf("SetInternalCodemappingDataHandler with %q - %q", input, err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Internal Codemapping Data saved successfully",
		"data":    input,
	})
}
