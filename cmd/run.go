package cmd

import (
	"fmt"
	"github.com/MaaXYZ/maa-framework-go"
	"github.com/MaaXYZ/maa-framework-go/toolkit"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run elf-aid-magic.",
	Run:   runRun,
}

func runRun(cmd *cobra.Command, args []string) {
	toolkit.InitOption("./", "{}")

	res := maa.NewResource(nil)
	defer res.Destroy()

	devices := toolkit.AdbDevices()
	if len(devices) == 0 {
		fmt.Println("No devices")
		return
	}

	fmt.Println("Devices:")
	for i, device := range devices {
		fmt.Printf("\t%d. %s(%s)\n", i, device.Address, device.Name)
		fmt.Printf("\t\t%s\n", device.AdbPath)
	}

	fmt.Println()
	fmt.Print("Select: ")
	var selectedDevice int
	_, err := fmt.Scanf("%d", &selectedDevice)
	if err != nil {
		return
	}
	device := devices[selectedDevice]

	ctrl := maa.NewAdbController(
		device.AdbPath,
		device.Address,
		device.ControllerType,
		device.Config,
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
		return
	}
}

func init() {
	rootCmd.AddCommand(runCmd)
}
