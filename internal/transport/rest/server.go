package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"medods-test-task/internal/config"
	"medods-test-task/internal/transport/rest/handler"
	"medods-test-task/internal/transport/rest/middleware"
)

type Server struct {
	router *fiber.App
	conf   config.Server
}

func New(conf config.Server) *Server {
	router := fiber.New(fiber.Config{
		ErrorHandler: middleware.Error(),
	})

	router.Use(
		recover.New(),
		logger.New(),
	)

	return &Server{
		router: router,
		conf:   conf,
	}
}

func (server *Server) Handle(sessionHandler *handler.SessionHandler) *Server {
	sessionHandler.Register(server.router)

	return server
}

func (server *Server) Listen() error {
	return server.router.Listen(server.conf.Addr)
}

func (server *Server) Shutdown() error {
	return server.router.Shutdown()
}
