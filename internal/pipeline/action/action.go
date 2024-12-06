package action

import (
	"github.com/MaaXYZ/maa-framework-go"
	"github.com/dongwlin/elf-aid-magic/internal/config"
	"github.com/dongwlin/elf-aid-magic/internal/gamemap"
	"go.uber.org/zap"
)

func Register(res *maa.Resource, conf *config.Config, logger *zap.Logger, taskerID string, navAsst *gamemap.NavigationAssistant) {
	res.RegisterCustomAction("SetCurrentLocation", NewSetCurrentLocationAction(logger, navAsst))
	res.RegisterCustomAction("MapNavigation", NewMapNavigationAction(logger, navAsst))
}
