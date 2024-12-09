package recognition

import (
	"encoding/json"
	"strconv"

	"github.com/MaaXYZ/maa-framework-go"
)

type UseRapidProjectileRecogniation struct{}

func NewUseRapidProjectileRecogniation() maa.CustomRecognition {
	return &UseRapidProjectileRecogniation{}
}

// Run implements maa.CustomRecognition.
func (a *UseRapidProjectileRecogniation) Run(ctx *maa.Context, arg *maa.CustomRecognitionArg) (*maa.CustomRecognitionResult, bool) {
	findRapidProjectileResult := ctx.RunRecognition("FindRapidProjectile", arg.Img)
	if findRapidProjectileResult == nil {
		return nil, false
	}
	if !findRapidProjectileResult.Hit {
		return nil, false
	}

	rapidProjectilesNumResult := ctx.RunRecognition("RapidProjectilesNum", arg.Img)
	if rapidProjectilesNumResult == nil {
		return nil, false
	}
	detailJson := rapidProjectilesNumResult.DetailJson
	var detail OCRDetail
	_ = json.Unmarshal([]byte(detailJson), &detail)
	num, err := strconv.Atoi(detail.Best.Text)
	if err != nil {
		return nil, false
	}
	if num == 0 {
		return nil, false
	}

	impactEscapeResponseResult := ctx.RunRecognition("ImpactEscapeResponse", arg.Img)
	if impactEscapeResponseResult != nil {
		if impactEscapeResponseResult.Hit {
			return nil, false
		}
	}

	return &maa.CustomRecognitionResult{
		Box: findRapidProjectileResult.Box,
	}, true
}
