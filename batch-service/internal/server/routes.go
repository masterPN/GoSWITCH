package server

import (
	"batch-service/internal/data"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	jsonContentType = "application/json"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	r.GET("/", s.HelloWorldHandler)
	r.POST("/internalCodemappingData", s.AddInternalCodemappingDataHandler)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) AddInternalCodemappingDataHandler(c *gin.Context) {
	var input data.InternalCodemappingData
	if err := c.BindJSON(&input); err != nil {
		c.Error(fmt.Errorf("AddInternalCodemappingDataHandler with %q - %q", input, err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := input.SendInternalCodemappingDataToRedis(); err != nil {
		c.Error(fmt.Errorf("AddInternalCodemappingDataHandler with %q - %q", input, err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "InternalCodemappingData added successfully",
		"data":    input,
	})
}
