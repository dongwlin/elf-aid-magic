package recognition

import (
	"encoding/json"

	"github.com/MaaXYZ/maa-framework-go"
)

type AutoAccelerationRecognition struct{}

func NewAutoAccelerationRecogniation() maa.CustomRecognition {
	return &AutoAccelerationRecognition{}
}

// Run implements maa.CustomRecognition.
func (a *AutoAccelerationRecognition) Run(ctx *maa.Context, arg *maa.CustomRecognitionArg) (*maa.CustomRecognitionResult, bool) {
	// recognizing rapid projectiles
	recRapidProjectiles := maa.J{
		"RapidProjectiles": maa.J{
			"recognition": TemplateMatch,
			"template":    "RapidProjectiles.png",
			"roi":         []int{1002, 599, 149, 121},
		},
	}
	ret := ctx.RunRecognition("RapidProjectiles", arg.Img, recRapidProjectiles)
	if ret == nil {
		return nil, false
	}
	if !ret.Hit {
		return nil, false
	}

	// check if the number of rapid projectiles is not 0
	recRapidProjectilesNum := maa.J{
		"RapidProjectilesNum": maa.J{
			"recognition": OCR,
			"roi":         []int{1059, 687, 36, 20},
			"replace":     []string{"。", "0"},
		},
	}
	ret = ctx.RunRecognition("RapidProjectilesNum", arg.Img, recRapidProjectilesNum)
	if ret == nil {
		return nil, false
	}
	detailJson := ret.DetailJson
	var detail OCRDetail
	_ = json.Unmarshal([]byte(detailJson), &detail)
	if detail.Best.Text == "0" {
		return nil, false
	}

	// check for imminent impact
	recStrike := maa.J{
		"Strike": maa.J{
			"recognition": OCR,
			"roi":         []int{973, 390, 73, 39},
			"expected":    []string{"请选择", "应对方式"},
		},
	}
	ret = ctx.RunRecognition("Strike", arg.Img, recStrike)
	if ret != nil {
		if ret.Hit {
			return nil, false
		}
	}

	return &maa.CustomRecognitionResult{
		Box: maa.Rect{X: 1002, Y: 599, W: 149, H: 121},
	}, true
}
