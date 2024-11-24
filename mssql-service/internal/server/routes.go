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
	r.GET("/operatorRouting", s.GetOperatorRoutingHandler)
	r.GET("/optimalRoute", s.ExecuteGetOptimalRouteHandler)
	r.GET("/internalCodemapping", s.FetchAllInternalCodemappingsHandler)
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
	if err := c.BindJSON(&input); err != nil {
		c.Error(fmt.Errorf("ExecuteRadiusOnestageValidateHandler with %q - %q", input, err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := s.onevoisModels.RadiusOnestageValidateData.ExecuteRadiusOnestageValidate(input.Prefix, input.CallingNumber, input.DestinationNumber)
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
	var input onevoisdata.RadiusAccounting
	if err := c.BindJSON(&input); err != nil {
		c.Error(fmt.Errorf("ExecuteRadiusAccountingHandler with %q - %q", input, err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

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

func (s *Server) GetOperatorRoutingHandler(c *gin.Context) {
	number := c.Query("number")

	result, err := s.wholesaleModels.ImgCdrOperatorRoutingData.GetFirstImgCdrOperatorRoutingByNumber(number)
	if err != nil {
		c.Error(fmt.Errorf("GetOperatorRoutingHandler with %q - %q", number, err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s *Server) ExecuteGetOptimalRouteHandler(c *gin.Context) {
	pCallString := c.Query("pCallString")

	result, err := s.wholesaleModels.OptimalRouteData.ExecuteGetOptimalRoute(pCallString)
	if err != nil {
		c.Error(fmt.Errorf("ExecuteGetOptimalRouteHandler with %q - %q", pCallString, err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s *Server) FetchAllInternalCodemappingsHandler(c *gin.Context) {
	result, err := s.wholesaleModels.InternalCodemappingData.GetAll()
	if err != nil {
		c.Error(fmt.Errorf("GetAllInternalCodemappingHandler - %q", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
