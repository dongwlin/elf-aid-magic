package handler

import "github.com/gofiber/fiber/v2"

type PingHandler struct{}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (h *PingHandler) Register(r fiber.Router) {
	r.All("/ping", h.Ping)
}

func (h *PingHandler) Ping(c *fiber.Ctx) error {
	return c.SendString("pong!!!!!")
}
