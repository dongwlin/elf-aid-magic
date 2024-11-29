package pipeline

import (
	"github.com/MaaXYZ/maa-framework-go"
	"github.com/dongwlin/elf-aid-magic/internal/config"
	"go.uber.org/zap"
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

func Init(res *maa.Resource, conf *config.Config, logger *zap.Logger) {
	res.RegisterCustomRecognition("RapidProjectiles", &AutoAccelerationRecognition{})
	res.RegisterCustomRecognition("IsAppInactive", NewIsAppInactiveRecognition(conf, logger))
}
