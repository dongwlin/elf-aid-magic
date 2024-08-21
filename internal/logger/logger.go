package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Level zapcore.Level

const (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
	PanicLevel = zapcore.PanicLevel
	FatalLevel = zapcore.FatalLevel
)

type Logger struct {
	*zap.Logger
}

type Config struct {
	Filename   string
	Level      Level
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	Dev        bool
}

func New(conf *Config) *Logger {
	checkedConf := check(conf)

	hook := lumberjack.Logger{
		Filename:   checkedConf.Filename,
		MaxSize:    checkedConf.MaxSize,
		MaxBackups: checkedConf.MaxBackups,
		MaxAge:     checkedConf.MaxAge,
		Compress:   checkedConf.Compress,
	}

	if checkedConf.Dev {
		encoder := getConsoleEncoder()
		core := zapcore.NewCore(
			encoder,
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)),
			zapcore.Level(conf.Level),
		)
		return &Logger{
			Logger: zap.New(core, zap.Development(), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)),
		}
	}
	encoder := getJsonEncoder()
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(&hook),
		zapcore.Level(conf.Level),
	)
	return &Logger{
		Logger: zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)),
	}
}

func check(conf *Config) Config {
	var checked Config
	if conf.Filename == "" {
		checked.Filename = "./eam.log"
	} else {
		checked.Filename = conf.Filename
	}
	checked.Level = conf.Level
	if conf.MaxSize == 0 {
		checked.MaxSize = 5
	} else {
		checked.MaxSize = conf.MaxSize
	}
	if conf.MaxBackups == 0 {
		checked.MaxBackups = 10
	} else {
		checked.MaxBackups = conf.MaxBackups
	}
	if conf.MaxAge == 0 {
		checked.MaxAge = 30
	} else {
		checked.MaxAge = conf.MaxAge
	}
	checked.Compress = conf.Compress

	return checked
}

func getConsoleEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
}

func getJsonEncoder() zapcore.Encoder {
	conf := zap.NewProductionEncoderConfig()
	conf.TimeKey = "time"
	conf.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewJSONEncoder(conf)
}
