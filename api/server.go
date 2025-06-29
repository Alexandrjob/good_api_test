package api

import (
	"github.com/gin-gonic/gin"
	"good_api_test/service"
)

type Server struct {
	service *service.Service
	router  *gin.Engine
}

func NewServer(service *service.Service) *Server {
	server := &Server{
		service,
		gin.Default(),
	}
	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	v1 := s.router.Group("/api/v1")
	{
		goodsGroup := v1.Group("/good")
		{
			goodsGroup.GET("", s.getGoodHandler)
			goodsGroup.POST("/create", s.createGoodHandler)
			goodsGroup.PATCH("/update", s.updateGoodHandler)
			goodsGroup.DELETE("/remove", s.deleteGoodHandler)
			goodsGroup.PATCH("/reprioritize", s.reprioritizeHandler)
		}

		v1.GET("/goods/list", s.getGoodsHandler)
	}
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}
