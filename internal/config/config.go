package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Device DeviceConfig `mapstructure:"device"`
	Tasks  []Task       `mapstructure:"tasks"`
	Log    LogConfig    `mapstructure:"log"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type DeviceConfig struct {
	Name         string                 `mapstructure:"name"`
	SerialNumber string                 `mapstructure:"serial_number"`
	AdbPath      string                 `mapstructure:"adb_path"`
	Input        uint64                 `mapstructure:"input"`
	Screencap    uint64                 `mapstructure:"screencap"`
	AdbConfig    map[string]interface{} `mapstructure:"adb_config"`
}

type Task struct {
	Entry string                 `json:"entry"`
	Param map[string]interface{} `json:"param"`
}

func New() *Config {
	v := viper.New()

	v.SetDefault("server.port", 8000)
	v.SetDefault("device.adb_config", map[string]interface{}{})

	v.SetConfigName("config")
	v.SetConfigType("toml")

	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path, %v", err)
	}
	exeDir := filepath.Dir(exePath)
	configDir := filepath.Join(exeDir, "config")
	v.AddConfigPath(configDir)

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file, %v", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		log.Fatalf("Failed to unmarshal config file, %v", err)
	}
	return &config
}
