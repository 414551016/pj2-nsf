package sbi

import (
	"net/http"
	"nsf/model"

	"github.com/gin-gonic/gin"
)

func (s *Server) getDiscRoutes() []route {
	return []route{
		{
			Name:        "Discovery Service",
			Method:      http.MethodGet,
			Pattern:     "",
			HandlerFunc: s.handleDiscoveryService,
		},
		{
			Name:        "Discovery Services",
			Method:      http.MethodGet,
			Pattern:     "/group",
			HandlerFunc: s.handleDiscoveryServices,
		},
	}
}

func (s *Server) handleDiscoveryService(c *gin.Context) {
	s.NsfLogger.SbiLog.Infoln("Handling Discovery Service")

	var request model.NnsfDiscGetRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		s.NsfLogger.SbiLog.Errorf("Failed to bind request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := s.Processor.ProcessDiscoveryService(request)
	if err != nil {
		s.NsfLogger.SbiLog.Errorln("Failed to process discovery service")
		c.JSON(err.StatusCode, gin.H{"error": err.ErrorDetail.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
	s.NsfLogger.SbiLog.Infoln("Discovery Service handled successfully")
}

func (s *Server) handleDiscoveryServices(c *gin.Context) {
	s.NsfLogger.SbiLog.Infoln("Handling Discovery Services")

	var request model.NnsfDiscGroupGetRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		s.NsfLogger.SbiLog.Errorln("Failed to bind request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := s.Processor.ProcessDiscoveryServices(request)
	if err != nil {
		s.NsfLogger.SbiLog.Errorln("Failed to process discovery services")
		c.JSON(err.StatusCode, gin.H{"error": err.ErrorDetail.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
	s.NsfLogger.SbiLog.Infoln("Discovery Services handled successfully")
}
