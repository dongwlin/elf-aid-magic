package handler

import (
	"github.com/dongwlin/elf-aid-magic/internal/logic"
	"github.com/gofiber/fiber/v2"
)

type PidHandler struct {
	pidLogic *logic.PidLogic
}

func NewPidHandler(pidLogic *logic.PidLogic) *PidHandler {
	return &PidHandler{
		pidLogic: logic.NewPidLogic(),
	}
}

func (h *PidHandler) Register(r fiber.Router) {
	r.Post("/pid/validate", h.ValidatePid)
}

type ValidatePidRequest struct {
	Pid int `json:"pid"`
}

func (h *PidHandler) ValidatePid(c *fiber.Ctx) error {
	var req ValidatePidRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"validated": false})
	}

	if req.Pid <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"validated": false})
	}

	validated := h.pidLogic.ValidatePid(req.Pid)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"validated": validated})
}
