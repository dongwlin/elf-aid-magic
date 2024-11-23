package cmd

import (
	"fmt"
	"os"

	"github.com/dongwlin/elf-aid-magic/internal/config"
	"github.com/dongwlin/elf-aid-magic/internal/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop eam server by daemon/pid file",
	Run:   stopRun,
}

func stopRun(cmd *cobra.Command, args []string) {
	conf := config.New()

	l := logger.New(conf)
	defer l.Sync()

	initDaemon(l)
	if pid == -1 {
		l.Warn("seems not have been started yet")
		fmt.Println("Seems not have been started. Try use `eam start` to start server.")
		return
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		l.Error("failed to find process", zap.Int("pid", pid), zap.Error(err))
		fmt.Printf("Failed to find process by pid: %d. See log.json for details.\n", pid)
		os.Exit(1)
	}
	err = process.Kill()
	if err != nil {
		l.Error("failed tod kill process", zap.Int("pid", pid), zap.Error(err))
		fmt.Printf("Failed to kill process by pid: %d. See log.json for details.\n", pid)
		os.Exit(1)
	} else {
		l.Info("killed process", zap.Int("pid", pid))
		fmt.Printf("Killed process by pid: %d\n", pid)
	}
	err = os.Remove(pidFile)
	if err != nil {
		l.Error("failed to remove pid file", zap.Error(err))
		fmt.Println("Failed to remove pid file. See log.json for details.")
	}
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
