package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/dongwlin/elf-aid-magic/internal/config"
	"github.com/dongwlin/elf-aid-magic/internal/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Silent start eam server",
	Run:   startRun,
}

func startRun(cmd *cobra.Command, args []string) {
	conf := config.New()

	l := logger.New(conf)
	defer l.Sync()

	initDaemon(l)
	if pid != -1 {
		_, err := os.FindProcess(pid)
		if err != nil {
			l.Info(
				"eam already started",
				zap.Int("pid", pid),
			)
			fmt.Printf("eam already started, pid: %d\n", pid)
			return
		}
	}

	serveArgs := os.Args
	serveArgs[1] = "serve"
	serve := &exec.Cmd{
		Path: serveArgs[0],
		Args: serveArgs,
		Env:  os.Environ(),
	}
	stdout, err := os.OpenFile(filepath.Join(filepath.Dir(pidFile), "start.log"), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		l.Error(
			"failed to open start file",
			zap.Error(err),
		)
		fmt.Println("Failed to open the start file. See log.json for details.")
		os.Exit(1)
	}
	serve.Stdout = stdout
	serve.Stderr = stdout
	if err = serve.Start(); err != nil {
		l.Error(
			"failed to start children process",
			zap.Error(err),
		)
		fmt.Println("Failed to start children process. See log.json for details.")
		os.Exit(1)
	}
	l.Info("success to start", zap.Int("pid", serve.Process.Pid))
	fmt.Printf("Success to start, pid: %d\n", serve.Process.Pid)
	err = os.WriteFile(pidFile, []byte(strconv.Itoa(serve.Process.Pid)), 0666)
	if err != nil {
		l.Warn(
			"failed to write pid file",
			zap.Error(err),
		)
		fmt.Println("Failed to record pid, you may not be able to stop the program with `eam stop`. See log.json for details.")
	}
}

func init() {
	rootCmd.AddCommand(startCmd)
}
