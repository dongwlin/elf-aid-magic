package cmd

import (
	"fmt"
	"github.com/dongwlin/elf-aid-magic/internal/logger"
	"os"
	"path/filepath"
)

func NewLogger() *logger.Logger {
	exe, err := os.Executable()
	if err != nil {
		fmt.Printf("Failed to get the path name for the executable, err: %v\n", err)
		os.Exit(1)
	}
	exeDir := filepath.Dir(exe)

	return logger.New(&logger.Config{
		Filename: filepath.Join(exeDir, "log", "log.json"),
		Dev:      dev,
	})
}
