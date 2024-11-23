package handler

import (
	"github.com/dongwlin/elf-aid-magic/internal/logic"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type VersionHandler struct {
	logger       *zap.Logger
	versionLogic *logic.VersionLogic
}

func NewVersionHandler(logger *zap.Logger, versionLogic *logic.VersionLogic) *VersionHandler {
	return &VersionHandler{
		logger:       logger,
		versionLogic: versionLogic,
	}
}

func (h *VersionHandler) Register(r fiber.Router) {
	r.Get("/versions/maa", h.GetMaaFrameworkVersion)
}

func (h *VersionHandler) GetMaaFrameworkVersion(c *fiber.Ctx) error {
	version := h.versionLogic.GetMaaFrameworkVersion()
	return c.Status(fiber.StatusOK).SendString(version)
}
