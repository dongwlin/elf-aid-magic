package pipeline

import (
	"github.com/MaaXYZ/maa-framework-go"
	"github.com/dongwlin/elf-aid-magic/internal/config"
	"github.com/dongwlin/elf-aid-magic/internal/gamemap"
	"github.com/dongwlin/elf-aid-magic/internal/pipeline/action"
	"github.com/dongwlin/elf-aid-magic/internal/pipeline/recognition"
	"go.uber.org/zap"
)

func Register(res *maa.Resource, conf *config.Config, logger *zap.Logger, taskerID string, navAsst *gamemap.NavigationAssistant) {
	action.Register(res, conf, logger, taskerID, navAsst)
	recognition.Register(res, conf, logger, taskerID)
}
