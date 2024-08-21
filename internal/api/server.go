package api

import "github.com/gofiber/fiber/v2"

type Server struct {
	router *fiber.App
}

func NewServer() *Server {
	router := fiber.New()

	return &Server{
		router: router,
	}
}

func (s *Server) Start(address string) error {
	return s.router.Listen(address)
}
