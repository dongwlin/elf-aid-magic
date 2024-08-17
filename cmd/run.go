package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/MaaXYZ/maa-framework-go"
	"github.com/MaaXYZ/maa-framework-go/toolkit"
	"github.com/dongwlin/elf-aid-magic/internal/config"
	"github.com/spf13/cobra"
	"os"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run elf-aid-magic.",
	Run:   runRun,
}

func runRun(cmd *cobra.Command, args []string) {
	toolkit.InitOption("./", "{}")

	conf, err := config.NewConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	res := maa.NewResource(nil)
	defer res.Destroy()

	var resJob maa.Job
	for _, resPath := range conf.Resource {
		resJob = res.PostPath(resPath)
	}
	resJob.Wait()

	adbConfigData, err := json.Marshal(conf.Adb.Config)
	if err != nil {
		fmt.Println(err)
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
	defer ctrl.Destroy()
	ctrlJob := ctrl.PostConnect()
	ctrlJob.Wait()

	inst := maa.New(nil)
	defer inst.Destroy()

	inst.BindResource(res)
	inst.BindController(ctrl)

	if !inst.Inited() {
		fmt.Println("Failed to initialize instance.")
		os.Exit(1)
	}

	var taskJob maa.TaskJob
	for _, task := range conf.Tasks {
		param, err := json.Marshal(task.Param)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		taskJob = inst.PostTask(task.Entry, string(param))
	}
	taskJob.Wait()
}

func init() {
	rootCmd.AddCommand(runCmd)
}
