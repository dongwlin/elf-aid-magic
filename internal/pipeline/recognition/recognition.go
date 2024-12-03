package recognition

import (
	"github.com/MaaXYZ/maa-framework-go"
	"github.com/dongwlin/elf-aid-magic/internal/config"
	"go.uber.org/zap"
)

const (
	DirectHit             = "DirectHit"
	TemplateMatch         = "TemplateMatch"
	FeatureMatch          = "FeatureMatch"
	ColorMatch            = "ColorMatch"
	OCR                   = "OCR"
	NeuralNetworkClassify = "NeuralNetworkClassify"
	NeuralNetworkDetect   = "NeuralNetworkDetect"
	Custom                = "Custom"
)

type OCRDetail struct {
	All      []OCRDetailItem `json:"all"`
	Best     OCRDetailItem   `json:"best"`
	Filtered []OCRDetailItem `json:"filtered"`
}

type OCRDetailItem struct {
	Box   []int   `json:"box"`
	Score float64 `json:"score"`
	Text  string  `json:"text"`
}

func Register(res *maa.Resource, conf *config.Config, logger *zap.Logger) {
	res.RegisterCustomRecognition("RapidProjectiles", NewAutoAccelerationRecogniation())
	res.RegisterCustomRecognition("IsAppInactive", NewIsAppInactiveRecognition(conf, logger))
}
