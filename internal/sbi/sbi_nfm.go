package sbi

import (
	"net/http"
	"nsf/model"

	"github.com/gin-gonic/gin"
)

// Network Function Management (NFM) Service
func (s *Server) getNfmRoutes() []route {
	return []route{
		{
			Name:        "Service Registration",
			Method:      http.MethodPost,
			Pattern:     "",
			HandlerFunc: s.handleServicesRegistration,
		},
		{
			Name:        "Services Deregistration",
			Method:      http.MethodDelete,
			Pattern:     "",
			HandlerFunc: s.handleServicesDeregistration,
		},
	}
}

func (s *Server) handleServicesRegistration(c *gin.Context) {
	s.NsfLogger.SbiLog.Infoln("Handling services registration")

	var request model.NnsfNfmPostRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		s.NsfLogger.SbiLog.Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := s.Processor.ProcessServicesRegistration(request)
	if err != nil {
		s.NsfLogger.SbiLog.Errorln("Failed to process services registration")
		c.JSON(err.StatusCode, gin.H{"error": err.ErrorDetail})
		return
	}

	c.JSON(http.StatusOK, response)
	s.NsfLogger.SbiLog.Infoln("Services registration handled successfully")
}

func (s *Server) handleServicesDeregistration(c *gin.Context) {
	s.NsfLogger.SbiLog.Infoln("Handling services deregistration")

	var request model.NnsfNfmDeleteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		s.NsfLogger.SbiLog.Errorln("Failed to bind request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.Processor.ProcessServicesDeregistration(request); err != nil {
		s.NsfLogger.SbiLog.Errorln("Failed to process services deregistration")
		c.JSON(err.StatusCode, gin.H{"error": err.ErrorDetail})
		return
	}

	c.JSON(http.StatusNoContent, nil)
	s.NsfLogger.SbiLog.Infoln("Services deregistration handled successfully")
}
