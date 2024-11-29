package operator

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/MaaXYZ/maa-framework-go"
	"github.com/dongwlin/elf-aid-magic/internal/config"
	"github.com/dongwlin/elf-aid-magic/internal/pipeline"
	"go.uber.org/zap"
)

type Operator struct {
	conf       *config.Config
	logger     *zap.Logger
	toolkit    *maa.Toolkit
	tasker     *maa.Tasker
	res        *maa.Resource
	ctrl       maa.Controller
	taskerID   string
	taskerName string
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
}

func (o *Operator) initToolkit() {
	toolkit := maa.NewToolkit()
	o.toolkit = toolkit
	o.toolkit.ConfigInitOption("./", "{}")
}

func (o *Operator) InitTasker() bool {
	return o.initTasker()
}

func (o *Operator) initTasker() bool {
	tasker := maa.NewTasker(nil)
	if tasker == nil {
		o.logger.Error("failed to init tasker.")
		return false
	}
	o.tasker = tasker
	return true
}

func (o *Operator) InitResource() bool {
	return o.initResource()
}

func (o *Operator) initResource() bool {
	res := maa.NewResource(nil)
	if res == nil {
		o.logger.Error("failed to init resource")
		return false
	}
	o.res = res
	exePath, err := os.Executable()
	if err != nil {
		o.logger.Error(
			"failed to get executable path",
			zap.Error(err),
		)
		return false
	}
	exeDir := filepath.Dir(exePath)
	resPath := filepath.Join(exeDir, "resource", "base")
	resJob := o.res.PostPath(resPath)
	o.logger.Info(
		"load resource",
		zap.String("resource path", resPath),
	)
	if ok := resJob.Wait().Success(); !ok {
		o.logger.Error(
			"failed to load resource",
			zap.String("resource", resPath),
		)
		return false
	}

	pipeline.Init(res, o.conf, o.logger)

	if ok := o.tasker.BindResource(o.res); !ok {
		o.logger.Error("failed to bind resource")
		return false
	}
	return true
}

func (o *Operator) InitController() bool {
	tasker, ok := o.getTaskerConfig()
	if !ok {
		return false
	}

	switch tasker.CtrlType {
	case "adb":
		return o.initAdbController()
	case "win32":
		return o.initWin32Controller()
	default:
		o.logger.Error(
			"unknown ctrl type",
			zap.String("ctrl type", tasker.CtrlType),
		)
		return false
	}
}

func (o *Operator) initAdbController() bool {
	tasker, ok := o.getTaskerConfig()
	if !ok {
		return false
	}
	device := tasker.AdbDevice

	var adbConfigStr string
	adbConfigData, err := json.Marshal(device.Config)
	if err != nil {
		o.logger.Error(
			"failed to serialize adb config",
			zap.Error(err),
		)
		return false
	}
	if device.Config == nil {
		adbConfigStr = "{}"
	} else {
		adbConfigStr = string(adbConfigData)
	}

	o.logger.Info(
		"adb config",
		zap.String("config", adbConfigStr),
	)

	ctrl := maa.NewAdbController(
		o.conf.AdbPath,
		device.SerialNumber,
		device.GetScreencapMethod(),
		device.GetInputMethod(),
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
		zap.String("path", o.conf.AdbPath),
		zap.String("address", device.SerialNumber),
	)
	if ok := o.tasker.BindController(o.ctrl); !ok {
		o.logger.Error("failed to bind controller")
		return false
	}
	return true
}

func (o *Operator) initWin32Controller() bool {
	tasker, ok := o.getTaskerConfig()
	if !ok {
		return false
	}
	window := tasker.Win32Window

	windows := o.toolkit.FindDesktopWindows()
	var handle unsafe.Pointer
	for _, window := range windows {
		if window.WindowName == "雷索纳斯" && window.ClassName == "UnityWndClass" {
			handle = window.Handle
			break
		}
	}
	if handle == nil {
		o.logger.Error("not found target window")
		return false
	}
	ctrl := maa.NewWin32Controller(
		handle,
		window.GetScreencapMethod(),
		window.GetInputMethod(),
		nil,
	)
	if ctrl == nil {
		o.logger.Error("failed to init win32 controller")
		return false
	}
	o.ctrl = ctrl
	o.logger.Info("create win32 controller")
	o.ctrl.SetScreenshotUseRawSize(true)
	if ok := o.tasker.BindController(o.ctrl); !ok {
		o.logger.Error("failed to bind controller")
		return false
	}
	return true
}

func (o *Operator) Connect() bool {
	if !o.ctrl.PostConnect().Wait().Success() {
		o.logger.Error("failed to connect")
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

	tasker, ok := o.getTaskerConfig()
	if !ok {
		return false
	}

	for _, task := range tasker.Tasks {
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

func (o *Operator) SetTaskerID(id string) {
	o.taskerID = id
}

func (o *Operator) SetTaskerName(name string) {
	o.taskerName = name
}

func (o *Operator) getTaskerConfig() (*config.TaskerConfig, bool) {
	taskers := o.conf.Taskers
	if len(taskers) == 0 {
		o.logger.Error("taskers is empty")
		return nil, false
	}

	logTaskerSelection := func(tasker *config.TaskerConfig) {
		o.logger.Info("selected tasker",
			zap.String("id", tasker.ID),
			zap.String("name", tasker.Name),
		)
	}

	if o.taskerID != "" && o.taskerName != "" {
		for _, tasker := range taskers {
			if tasker.ID == o.taskerID && tasker.Name == o.taskerName {
				logTaskerSelection(tasker)
				return tasker, true
			}
		}
		o.logger.Error("no tasker found with specified id and name",
			zap.String("id", o.taskerID),
			zap.String("name", o.taskerName),
		)
		return nil, false
	}

	if o.taskerID != "" {
		for _, tasker := range taskers {
			if tasker.ID == o.taskerID {
				logTaskerSelection(tasker)
				return tasker, true
			}
		}
		o.logger.Error("no tasker found with specified id",
			zap.String("id", o.taskerID),
		)
		return nil, false
	}

	if o.taskerName != "" {
		for _, tasker := range taskers {
			if tasker.Name == o.taskerName {
				logTaskerSelection(tasker)
				return tasker, true
			}
		}
		o.logger.Warn("no tasker found with specified name",
			zap.String("name", o.taskerName),
		)
		return nil, false
	}

	o.logger.Info("no specific tasker id or name provided, defaulting to the first tasker")
	logTaskerSelection(taskers[0])
	return taskers[0], true
}
