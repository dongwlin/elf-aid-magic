package logic

import (
	"github.com/MaaXYZ/maa-framework-go"
	"github.com/dongwlin/elf-aid-magic/internal/config"
)

type VersionLogic struct{}

func NewVersionLogic() *VersionLogic {
	return &VersionLogic{}
}

func (l *VersionLogic) GetMaaFrameworkVersion() string {
	return maa.Version()
}

func (l *VersionLogic) GetElfAidMagicVersion() string {
	return config.Version
}

func (l *VersionLogic) GetGoVersion() string {
	return config.GoVersion
}

func (l *VersionLogic) GetBuildAt() string {
	return config.BuildAt
}
