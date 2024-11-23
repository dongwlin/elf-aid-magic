package cmd

import (
	"fmt"

	"github.com/dongwlin/elf-aid-magic/internal/logic"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of elf-aid-magic.",
	Run:   versionRun,
}

func versionRun(cmd *cobra.Command, args []string) {
	versionLogic := logic.NewVersionLogic()
	fmt.Println("Build At:", versionLogic.GetBuildAt())
	fmt.Println("Go Version:", versionLogic.GetGoVersion())
	fmt.Println("Version:", versionLogic.GetElfAidMagicVersion())
	fmt.Println("Maa Framework Version:", versionLogic.GetMaaFrameworkVersion())
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
