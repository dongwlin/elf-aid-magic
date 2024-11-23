package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"go.uber.org/zap"
)

var (
	pid     = -1
	pidFile string
)

func initDaemon(logger *zap.Logger) {
	exe, err := os.Executable()
	if err != nil {
		logger.Error("failed to get the path name for the executable", zap.Error(err))
		fmt.Println("Failed to get the path name for the executable. See log.json for details.")
		os.Exit(1)
	}

	exeDir := filepath.Dir(exe)
	_ = os.MkdirAll(filepath.Join(exeDir, "daemon"), 0700)
	pidFile = filepath.Join(exeDir, "daemon", "pid")
	if _, err = os.Stat(pidFile); err == nil {
		pidData, err := os.ReadFile(pidFile)
		if err != nil {
			logger.Error("failed to read the pid file", zap.Error(err))
			fmt.Println("Failed to read pid file. See log.json for details.")
			os.Exit(1)
		}
		pid, err = strconv.Atoi(string(pidData))
		if err != nil {
			logger.Error("failed to parse the pid file", zap.Error(err))
			fmt.Println("Failed to parse pid file. See log.json for details.")
			os.Exit(1)
		}
	}
}
