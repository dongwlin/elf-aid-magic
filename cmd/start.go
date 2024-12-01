package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/dongwlin/elf-aid-magic/internal/config"
	"github.com/dongwlin/elf-aid-magic/internal/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Silent start eam server",
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

func start() {
	conf := config.New()

	l := logger.New(conf)
	defer l.Sync()

	initDaemon(l)
	if pid != -1 && isProcessRunning(conf, pid) {
		l.Info("eam already started",
			zap.Int("pid", pid),
		)
		fmt.Printf("eam already started, pid: %d\n", pid)
		return
	}
	if isServerRunning(conf) {
		l.Warn("eam may be running, but not the current PID",
			zap.Int("pid", pid),
		)
		fmt.Printf("eam may be running, but not the current PID: %d\n", pid)
		return
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
		l.Error("failed to open start file",
			zap.Error(err),
		)
		fmt.Println("Failed to open the start file. See log.jsonl for details.")
		os.Exit(1)
	}
	serve.Stdout = stdout
	serve.Stderr = stdout
	if err = serve.Start(); err != nil {
		l.Error("failed to start children process",
			zap.Error(err),
		)
		fmt.Println("Failed to start children process. See log.jsonl for details.")
		os.Exit(1)
	}

	// Wait for a short period to allow the server to start
	time.Sleep(5 * time.Second)

	if !isServerRunning(conf) {
		l.Error("failed to verify server start")
		fmt.Println("Failed to verify server start. See log.jsonl for details.")
		os.Exit(1)
	}

	l.Info("success to start",
		zap.Int("pid", serve.Process.Pid),
	)
	fmt.Printf("Success to start, pid: %d\n", serve.Process.Pid)
	fmt.Printf("Local: http://localhost:%d\n", conf.Server.Port)

	err = os.WriteFile(pidFile, []byte(strconv.Itoa(serve.Process.Pid)), 0666)
	if err != nil {
		l.Warn("failed to write pid file",
			zap.Error(err),
		)
		fmt.Println("Failed to record pid, you may not be able to stop the program with `eam stop`. See log.jsonl for details.")
	}
}

func isProcessRunning(conf *config.Config, pid int) bool {
	reqBody, err := json.Marshal(map[string]int{"pid": pid})
	if err != nil {
		return false
	}
	url := fmt.Sprintf("http://localhost:%d/pid/validate", conf.Server.Port)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	var result map[string]bool
	if err := json.Unmarshal(body, &result); err != nil {
		return false
	}

	return resp.StatusCode == http.StatusOK && result["validated"]
}

func isServerRunning(conf *config.Config) bool {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/ping", conf.Server.Port))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	return resp.StatusCode == http.StatusOK && string(body) == "pong!!!!!"
}

func init() {
	rootCmd.AddCommand(startCmd)
}
