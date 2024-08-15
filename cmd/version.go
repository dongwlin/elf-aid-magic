package cmd

import (
	"fmt"
	"github.com/MaaXYZ/maa-framework-go"
	"github.com/dongwlin/elf-aid-magic/internal/config"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of elf-aid-magic.",
	Run:   versionRun,
}

func versionRun(cmd *cobra.Command, args []string) {
	maaVersion := maa.Version()
	fmt.Println("Build At:", config.BuildAt)
	fmt.Println("Go Version:", config.GoVersion)
	fmt.Println("Version:", config.Version)
	fmt.Println("Maa Framework Version:", maaVersion)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
