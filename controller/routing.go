package controller

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/skhaz/scheduler/repository"
	"github.com/skhaz/scheduler/workflow"
	"go.uber.org/zap"
)

type Server struct {
	router *gin.Engine
}

func InitServer() *Server {
	return &Server{router: gin.Default()}
}

func (s *Server) Run() {
	s.registerRoutes()

	s.router.Use(gzip.Gzip(gzip.DefaultCompression))

	_ = s.router.Run()
}

func (s *Server) SetLogger(logger *zap.Logger) {
	s.router.Use(func(c *gin.Context) {
		c.Set("Logger", logger)
		c.Next()
	})

	// s.router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	// s.router.Use(ginzap.RecoveryWithZap(logger, true))
}

func (s *Server) SetRepositoryRegistry(rr *repository.RepositoryRegistry) {
	s.router.Use(func(c *gin.Context) {
		c.Set("RepositoryRegistry", rr)
		c.Next()
	})
}

func (s *Server) SetWorkflow(wf workflow.Interface) {
	s.router.Use(func(c *gin.Context) {
		c.Set("Workflow", wf)
		c.Next()
	})
}

func (s *Server) registerRoutes() {
	var router = s.router

	router.NoRoute(NoRoute)

	triggers := router.Group("/triggers")
	{
		triggers.GET("", GetTriggers)
		triggers.POST("", CreateTrigger)
		triggers.GET("/:uuid", GetTrigger)
		triggers.DELETE("/:uuid", DeleteTrigger)
	}
}
