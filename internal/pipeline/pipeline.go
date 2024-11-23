package pipeline

import "github.com/MaaXYZ/maa-framework-go"

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

func Init(res *maa.Resource) {
	res.RegisterCustomRecognition("RapidProjectiles", &AutoAccelerationRecognition{})
}
