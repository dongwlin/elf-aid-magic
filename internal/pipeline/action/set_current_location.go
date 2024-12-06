package action

import (
	"encoding/json"

	"github.com/MaaXYZ/maa-framework-go"
	"github.com/dongwlin/elf-aid-magic/internal/gamemap"
	"github.com/dongwlin/elf-aid-magic/internal/pipeline/recognition"
	"go.uber.org/zap"
)

type SetCurrentLocationAction struct {
	logger  *zap.Logger
	navAsst *gamemap.NavigationAssistant
}

func NewSetCurrentLocationAction(logger *zap.Logger, navAsst *gamemap.NavigationAssistant) maa.CustomAction {
	return &SetCurrentLocationAction{
		logger:  logger,
		navAsst: navAsst,
	}
}

// Run implements maa.CustomAction.
func (a *SetCurrentLocationAction) Run(ctx *maa.Context, arg *maa.CustomActionArg) bool {
	a.logger.Info("Starting SetCurrentLocationAction")

	ctrl := ctx.GetTasker().GetController()
	ctrl.PostScreencap().Wait()
	a.logger.Debug("Screencap posted")

	img := ctrl.CacheImage()
	a.logger.Debug("Image cached for recognition")

	CurrentLocationResult := ctx.RunRecognition("CurrentLocation", img)
	if CurrentLocationResult == nil {
		a.logger.Error("Recognition result is nil")
		return false
	}
	detailJson := CurrentLocationResult.DetailJson
	var detail recognition.OCRDetail
	err := json.Unmarshal([]byte(detailJson), &detail)
	if err != nil {
		a.logger.Error("Failed to unmarshal recognition detail JSON",
			zap.Error(err),
		)
		return false
	}

	var name string
	for _, item := range detail.Filtered {
		var exists bool
		name, exists = gamemap.GetLocationNameByZhCN(item.Text)
		if exists {
			break
		}
	}
	if name == "" {
		a.logger.Error("Failed to get location name from recognized text")
		return false
	}

	location, exists := gamemap.GetLocation(name)
	if !exists {
		a.logger.Warn("Location does not exist",
			zap.String("name", name),
		)
		return false
	}

	a.navAsst.SetCurrentLocation(location)
	a.logger.Info("Current location set successfully",
		zap.String("location", name),
	)

	return true
}
