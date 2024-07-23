package server

import (
	"mssql-service/internal/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)
	r.POST("/callLog1801", s.CreateCallLog1801Handler)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) CreateCallLog1801Handler(c *gin.Context) {
	var input data.CallLog
	c.BindJSON(&input)

	// Call the Insert() on cdrs model
	err := s.models.CallLogs.Insert(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		panic(err)
	}

	resp := make(map[string]string)
	resp["message"] = "success"

	c.JSON(http.StatusOK, resp)
}
