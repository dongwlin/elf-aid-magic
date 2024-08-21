package api

import "github.com/gofiber/fiber/v2"

func (s *Server) ping(c *fiber.Ctx) error {
	return c.SendString("pong!!!!!")
}
