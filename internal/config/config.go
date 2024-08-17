package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Adb      AdbController `json:"adb"`
	Resource []string      `json:"resource"`
	Tasks    []Task        `json:"tasks"`
}

type AdbController struct {
	Path      string                 `json:"path"`
	Address   string                 `json:"address"`
	Key       int32                  `json:"key"`
	Touch     int32                  `json:"touch"`
	Screencap int32                  `json:"screencap"`
	Config    map[string]interface{} `json:"config"`
}

type Task struct {
	Entry string                 `json:"entry"`
	Param map[string]interface{} `json:"param"`
}

func NewConfig() (*Config, error) {
	data, err := load("./config/eam_config.json")
	if err != nil {
		return nil, err
	}

	conf := &Config{}
	err = json.Unmarshal(data, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func load(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}
