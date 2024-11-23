package operator

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/MaaXYZ/maa-framework-go"
	"github.com/dongwlin/elf-aid-magic/internal/config"
	"github.com/dongwlin/elf-aid-magic/internal/pipeline"
	"go.uber.org/zap"
)

type Operator struct {
	conf    *config.Config
	logger  *zap.Logger
	toolkit *maa.Toolkit
	tasker  *maa.Tasker
	res     *maa.Resource
	ctrl    maa.Controller
}

func New(conf *config.Config, logger *zap.Logger) *Operator {
	o := &Operator{
		conf:   conf,
		logger: logger,
	}
	o.init()
	return o
}

func (o *Operator) Destroy() {
	if o.ctrl != nil {
		o.ctrl.Destroy()
	}
	if o.res != nil {
		o.res.Destroy()
	}
	if o.tasker != nil {
		o.tasker.Destroy()
	}
}

func (o *Operator) init() {
	o.initToolkit()
	o.initTasker()
	o.initResource()
}

func (o *Operator) initToolkit() {
	toolkit := maa.NewToolkit()
	o.toolkit = toolkit
	o.toolkit.ConfigInitOption("./", "{}")
}

func (o *Operator) initTasker() {
	tasker := maa.NewTasker(nil)
	o.tasker = tasker
}

func (o *Operator) initResource() {
	res := maa.NewResource(nil)
	o.res = res
	exePath, err := os.Executable()
	if err != nil {
		o.Destroy()
		o.logger.Fatal(
			"failed to get executable path",
			zap.Error(err),
		)
	}
	exeDir := filepath.Dir(exePath)
	resPath := filepath.Join(exeDir, "resource", "base")
	resJob := o.res.PostPath(resPath)
	o.logger.Info(
		"load resource",
		zap.String("resource path", resPath),
	)
	if ok := resJob.Wait().Success(); !ok {
		o.Destroy()
		o.logger.Fatal(
			"failed to load resource",
			zap.String("resource", resPath),
		)
	}

	pipeline.Init(res)

	if ok := o.tasker.BindResource(o.res); !ok {
		o.Destroy()
		o.logger.Fatal("failed to bind resource")
	}
}

func (o *Operator) initController() bool {
	var adbConfigStr string
	adbConfigData, err := json.Marshal(o.conf.Device.AdbConfig)
	if err != nil {
		o.logger.Error(
			"failed to serialize adb config",
			zap.Error(err),
		)
		return false
	}
	if o.conf.Device.AdbConfig == nil {
		adbConfigStr = "{}"
	} else {
		adbConfigStr = string(adbConfigData)
	}

	o.logger.Info(
		"adb config",
		zap.String("config", adbConfigStr),
	)

	ctrl := maa.NewAdbController(
		o.conf.Device.AdbPath,
		o.conf.Device.SerialNumber,
		maa.AdbScreencapMethod(o.conf.Device.Screencap),
		maa.AdbInputMethod(o.conf.Device.Input),
		adbConfigStr,
		"./MaaAgentBinary",
		nil,
	)
	if ctrl == nil {
		o.logger.Error("failed to init adb controller")
		return false
	}
	o.ctrl = ctrl
	o.logger.Info(
		"create adb controller",
		zap.String("path", o.conf.Device.AdbPath),
		zap.String("address", o.conf.Device.SerialNumber),
	)
	if ok := ctrl.PostConnect().Wait().Success(); !ok {
		o.logger.Error(
			"failed to connect device",
			zap.String("path", o.conf.Device.AdbPath),
			zap.String("address", o.conf.Device.SerialNumber),
		)
		return false
	}
	if ok := o.tasker.BindController(ctrl); !ok {
		o.logger.Error("failed to bind controller")
		return false
	}
	return true
}

func (o *Operator) Connect() bool {
	if !o.initController() {
		return false
	}
	if !o.tasker.Initialized() {
		o.logger.Error("failed to initialize tasker instance")
		return false
	}
	return true
}

func (o *Operator) Run(ctx context.Context) bool {
	if !o.tasker.Initialized() {
		o.logger.Error("failed to initialize tasker instance")
		return false
	}

	for _, task := range o.conf.Tasks {
		select {
		case <-ctx.Done():
			o.logger.Info("operation cancelled")
			return false
		default:
		}

		param, err := json.Marshal(task.Param)
		if err != nil {
			o.Destroy()
			o.logger.Fatal(
				"failed to serialize task param",
				zap.Error(err),
			)
		}
		o.logger.Info(
			"run task",
			zap.String("entry", task.Entry),
			zap.String("param", string(param)),
		)
		if ok := o.tasker.PostPipeline(task.Entry, string(param)).Wait().Success(); !ok {
			o.logger.Error(
				"failed to complete the task",
				zap.String("entry", task.Entry),
			)
			select {
			case <-ctx.Done():
				o.logger.Info("operation cancelled")
				return false
			default:
				continue
			}
		}
		o.logger.Info(
			"success to complete the task",
			zap.String("entry", task.Entry),
		)
	}
	o.logger.Info("complete all tasks")
	return true
}

func (o *Operator) Stop() bool {
	return o.tasker.PostStop()
}
