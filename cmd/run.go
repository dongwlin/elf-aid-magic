package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/MaaXYZ/maa-framework-go"
	"github.com/MaaXYZ/maa-framework-go/toolkit"
	"github.com/dongwlin/elf-aid-magic/internal/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run elf-aid-magic.",
	Run:   runRun,
}

func runRun(_ *cobra.Command, _ []string) {
	l := NewLogger()
	defer l.Sync()

	toolkit.InitOption("./", "{}")

	conf, err := config.NewConfig()
	if err != nil {
		l.Error("fail to init config", zap.Error(err))
		fmt.Println("Fail to init config. See log.json for details.")
		os.Exit(1)
	}

	inst := maa.New(nil)
	defer inst.Destroy()

	res := maa.NewResource(nil)
	defer res.Destroy()

	for _, resPath := range conf.Resource {
		resJob := res.PostPath(resPath)
		l.Info("load resource", zap.String("resource", resPath))
		if ok := resJob.Wait(); !ok {
			l.Error("fail to load resource", zap.String("resource", resPath))
			fmt.Println("Fail to load resource. See log.json for details.")
			os.Exit(1)
		}
	}
	if ok := inst.BindResource(res); !ok {
		l.Error("failed to bind resource")
		fmt.Println("Fail to bind resource. See log.json for details.")
		os.Exit(1)
	}

	adbConfigData, err := json.Marshal(conf.Adb.Config)
	if err != nil {
		l.Error("failed to serialize adb config", zap.Error(err))
		fmt.Println("Failed to serialize adb config. See log.json for details.")
		os.Exit(1)
	}

	ctrl := maa.NewAdbController(
		conf.Adb.Path,
		conf.Adb.Address,
		maa.AdbControllerType(conf.Adb.Key|conf.Adb.Touch|conf.Adb.Screencap),
		string(adbConfigData),
		"./MaaAgentBinary",
		nil,
	)
	if ctrl == nil {
		l.Error("failed to init adb controller")
		fmt.Println("Failed to init adb controller. See log.json for details.")
		os.Exit(1)
	}
	defer ctrl.Destroy()
	l.Info("new adb controller", zap.String("path", conf.Adb.Path), zap.String("address", conf.Adb.Address))
	if ok := ctrl.PostConnect().Wait(); !ok {
		l.Error("failed to connect device", zap.String("path", conf.Adb.Path), zap.String("address", conf.Adb.Address))
		fmt.Println("Failed to connect device. See log.json for details.")
		os.Exit(1)
	}
	if ok := inst.BindController(ctrl); !ok {
		l.Error("failed to bind controller")
		fmt.Println("Fail to bind controller. See log.json for details.")
		os.Exit(1)
	}

	if !inst.Inited() {
		l.Error("failed to initialize instance")
		fmt.Println("Failed to initialize instance. See log.json for details.")
		os.Exit(1)
	}

	for _, task := range conf.Tasks {
		param, err := json.Marshal(task.Param)
		if err != nil {
			l.Error("failed to serialize task param", zap.Error(err))
			fmt.Println("Failed to serialize task param. See log.json for details.")
			os.Exit(1)
		}
		l.Info("run task", zap.String("entry", task.Entry), zap.String("param", string(param)))
		if ok := inst.PostTask(task.Entry, string(param)).Wait(); !ok {
			l.Error("failed to complete task", zap.String("entry", task.Entry))
		}
		l.Info("success to complete task", zap.String("entry", task.Entry))
	}
	l.Info("complete all tasks")
}

func init() {
	rootCmd.AddCommand(runCmd)
}
