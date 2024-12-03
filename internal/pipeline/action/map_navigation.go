package action

import (
	"fmt"
	"math"
	"time"

	"github.com/MaaXYZ/maa-framework-go"
	"github.com/dongwlin/elf-aid-magic/internal/gamemap"
	"go.uber.org/zap"
)

type MapNavigationAction struct {
	logger  *zap.Logger
	navAsst *gamemap.NavigationAssistant
}

func NewMapNavigationAction(logger *zap.Logger, navAsst *gamemap.NavigationAssistant) maa.CustomAction {
	return &MapNavigationAction{
		logger:  logger,
		navAsst: navAsst,
	}
}

// Run implements maa.CustomAction.
func (a *MapNavigationAction) Run(ctx *maa.Context, arg *maa.CustomActionArg) bool {
	destName := gamemap.CapeCity
	a.findDestination(ctx, destName)

	ctrl := ctx.GetTasker().GetController()
	dest, exists := gamemap.GetLocation(destName)
	if !exists {
		a.logger.Error("dest not exists",
			zap.String("loaction", destName),
		)
	}

	angle := a.navAsst.AngleTo(dest)
	center := gamemap.Point{
		X: 640,
		Y: 360,
	}
	start, end := a.getCirclePoints(center, 280, angle)

	for !a.findDestination(ctx, destName) {
		ctrl.PostSwipe(int32(start.X), int32(start.Y), int32(end.X), int32(end.Y), 500*time.Millisecond)
	}

	return true
}

func (a *MapNavigationAction) findDestination(ctx *maa.Context, dest string) bool {
	ctrl := ctx.GetTasker().GetController()
	ctrl.PostScreencap().Wait()
	img := ctrl.CacheImage()
	task := fmt.Sprintf("To%s", dest)
	ret := ctx.RunRecognition(task, img)
	if ret == nil {
		return false
	}
	if ret.Hit {
		ctx.RunAction(task, ret.Box, ret.DetailJson)
		return true
	}
	return false
}

// getCirclePoints calculates the two points on the circle for a given angle.
func (a *MapNavigationAction) getCirclePoints(center gamemap.Point, radius float64, angle float64) (gamemap.Point, gamemap.Point) {
	// Convert angle from degrees to radians
	radians := angle * (math.Pi / 180.0)

	centerX := float64(center.X)
	centerY := float64(center.Y)

	// Calculate the first point
	x1 := centerX + radius*math.Cos(radians)
	y1 := centerY + radius*math.Sin(radians)

	// Calculate the opposite angle (180 degrees apart)
	oppositeRadians := radians + math.Pi

	// Calculate the second point
	x2 := centerX + radius*math.Cos(oppositeRadians)
	y2 := centerY + radius*math.Sin(oppositeRadians)

	p1 := gamemap.Point{
		X: int(x1),
		Y: int(y1),
	}
	p2 := gamemap.Point{
		X: int(x2),
		Y: int(y2),
	}

	return p1, p2
}
