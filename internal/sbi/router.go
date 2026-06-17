package sbi

import "github.com/gin-gonic/gin"

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc gin.HandlerFunc
}

func addService(s *Server) {
	s.router.RedirectTrailingSlash = false

	sbiGroup := s.router.Group("/sbi")

	addRoutes(sbiGroup.Group("/heartbeat"), s.getHeartbeat())
	addRoutes(sbiGroup.Group("/nnsf-nfm"), s.getNfmRoutes())
	addRoutes(sbiGroup.Group("/nnsf-disc"), s.getDiscRoutes())
}

func addRoutes(group *gin.RouterGroup, routes []route) {
	for _, route := range routes {
		switch route.Method {
		case "GET":
			group.GET(route.Pattern, route.HandlerFunc)
		case "POST":
			group.POST(route.Pattern, route.HandlerFunc)
		case "PUT":
			group.PUT(route.Pattern, route.HandlerFunc)
		case "DELETE":
			group.DELETE(route.Pattern, route.HandlerFunc)
		case "PATCH":
			group.PATCH(route.Pattern, route.HandlerFunc)
		}
	}
}
