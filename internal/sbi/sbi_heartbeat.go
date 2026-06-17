package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getHeartbeat() []route {
	return []route{
		{
			Name:        "Heartbeat",
			Method:      http.MethodGet,
			Pattern:     "",
			HandlerFunc: s.handleHeartbeat,
		},
	}
}

func (s *Server) handleHeartbeat(c *gin.Context) {
	s.NsfLogger.SbiLog.Infoln("Handling heartbeat")

	c.Status(http.StatusNoContent)

	s.NsfLogger.SbiLog.Infoln("Heartbeat handled successfully")
}
