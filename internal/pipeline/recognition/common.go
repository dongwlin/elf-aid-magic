package recognition

import (
	"encoding/json"

	"github.com/MaaXYZ/maa-framework-go"
	"github.com/dongwlin/elf-aid-magic/internal/config"
	"github.com/dongwlin/elf-aid-magic/internal/pkg/adbtool"
	"go.uber.org/zap"
)

type IsAppInactiveRecognition struct {
	conf   *config.Config
	logger *zap.Logger
}

func NewIsAppInactiveRecognition(conf *config.Config, logger *zap.Logger) maa.CustomRecognition {
	return &IsAppInactiveRecognition{
		conf:   conf,
		logger: logger,
	}
}

type IsAppInactiveRecognitionRunParam struct {
	Package string `json:"package"`
}

// Run implements maa.CustomRecognition.
func (i *IsAppInactiveRecognition) Run(ctx *maa.Context, arg *maa.CustomRecognitionArg) (*maa.CustomRecognitionResult, bool) {
	taskers := i.conf.Taskers
	if len(taskers) == 0 {
		i.logger.Error("taskers is empty")
		return nil, false
	}

	tasker := taskers[0]
	device := tasker.AdbDevice

	var param IsAppInactiveRecognitionRunParam
	err := json.Unmarshal([]byte(arg.CustomRecognitionParam), &param)
	if err != nil {
		return nil, false
	}
	actived, err := adbtool.IsAppActive(i.conf.AdbPath, device.SerialNumber, param.Package)
	if err != nil {
		return nil, false
	}
	if actived {
		return nil, false
	}
	return &maa.CustomRecognitionResult{
		Box: maa.Rect{X: 0, Y: 0, W: 0, H: 0},
	}, true
}
