package server

import (
	"fmt"
	"mssql-service/internal/data/onevoisdata"
	"mssql-service/internal/data/wholesaledata"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	routeHelloWorld             = "/"
	routeOperatorRouting        = "/operatorRouting"
	routeOptimalRoute           = "/optimalRoute"
	routeRadiusOnestageValidate = "/radiusOnestageValidate"
	routeRadiusAccounting       = "/radiusAccounting"
	routeInternalCodemapping    = "/internalCodemapping"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET(routeHelloWorld, s.HelloWorldHandler)
	r.GET(routeOperatorRouting, s.GetOperatorRoutingHandler)
	r.GET(routeOptimalRoute, s.ExecuteGetOptimalRouteHandler)
	r.POST(routeRadiusOnestageValidate, s.ExecuteRadiusOnestageValidateHandler)
	r.POST(routeRadiusAccounting, s.ExecuteRadiusAccountingHandler)
	r.GET(routeInternalCodemapping, s.FetchAllInternalCodemappingsHandler)
	r.POST(routeInternalCodemapping, s.SetInternalCodemappingHandler)
	r.DELETE(routeInternalCodemapping, s.DeleteInternalCodemappingHandler)

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
	if err := c.ShouldBindJSON(&input); err != nil {
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
	if err := c.ShouldBindJSON(&input); err != nil {
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

func (s *Server) SetInternalCodemappingHandler(c *gin.Context) {
	var input wholesaledata.InternalCodemappingData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(fmt.Errorf("SetInternalCodemappingHandler with %q - %q", input, err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := s.wholesaleModels.InternalCodemappingData.Set(input)
	if err != nil {
		c.Error(fmt.Errorf("SetInternalCodemappingHandler with %q - %q", input, err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s *Server) DeleteInternalCodemappingHandler(c *gin.Context) {
	internalCodeString := c.Query("internalCode")
	internalCode, err := strconv.Atoi(internalCodeString)
	if err != nil {
		c.Error(fmt.Errorf("DeleteInternalCodemappingHandler with id %q - %q", internalCodeString, err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = s.wholesaleModels.InternalCodemappingData.Delete(internalCode)
	if err != nil {
		c.Error(fmt.Errorf("DeleteInternalCodemappingHandler with id %q - %q", internalCodeString, err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Internal codemapping deleted successfully",
		"internalCode": internalCode,
	})
}
