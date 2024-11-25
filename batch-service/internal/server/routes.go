package server

import (
	"batch-service/internal/data"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
		s.handleError(c, http.StatusBadRequest, err)
		return
	}

	mssqlResp, mssqlErr := input.SendInternalCodemappingDataToMssql()
	if mssqlErr != nil {
		s.handleError(c, http.StatusInternalServerError, fmt.Errorf("failed to send data to MSSQL: %s", mssqlErr.Error()))
		return
	}

	redisErr := mssqlResp.SendInternalCodemappingDataToRedis()
	if redisErr != nil {
		s.handleError(c, http.StatusInternalServerError, fmt.Errorf("failed to send data to Redis: %s", redisErr.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "InternalCodemappingData added successfully",
		"data":    mssqlResp,
	})
}

func (s *Server) handleError(c *gin.Context, statusCode int, err error) {
	c.Error(err)
	c.JSON(statusCode, gin.H{
		"error": err.Error(),
	})
}
