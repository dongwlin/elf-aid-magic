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
	version := r.Group("/versions")
	version.Get("/", h.GetVersions)
	version.Get("/maa", h.GetMaaFrameworkVersion)
	version.Get("/eam", h.GetElfAidMagicVersion)
	version.Get("/go", h.GetGoVersion)
}

func (h *VersionHandler) GetVersions(c *fiber.Ctx) error {
	maaVersion := h.versionLogic.GetMaaFrameworkVersion()
	eamVersion := h.versionLogic.GetElfAidMagicVersion()
	goVersion := h.versionLogic.GetGoVersion()
	buildAt := h.versionLogic.GetBuildAt()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"maa_version": maaVersion,
		"eam_version": eamVersion,
		"go_version":  goVersion,
		"build_at":    buildAt,
	})
}

func (h *VersionHandler) GetMaaFrameworkVersion(c *fiber.Ctx) error {
	version := h.versionLogic.GetMaaFrameworkVersion()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"version": version,
	})
}

func (h *VersionHandler) GetElfAidMagicVersion(c *fiber.Ctx) error {
	version := h.versionLogic.GetElfAidMagicVersion()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"version": version,
	})
}

func (h *VersionHandler) GetGoVersion(c *fiber.Ctx) error {
	version := h.versionLogic.GetGoVersion()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"version": version,
	})
}

func (h *VersionHandler) GetBuildAt(c *fiber.Ctx) error {
	buildAt := h.versionLogic.GetBuildAt()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"build_at": buildAt,
	})
}
