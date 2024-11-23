package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dongwlin/elf-aid-magic/internal/config"
	"github.com/dongwlin/elf-aid-magic/internal/logger"
	"github.com/dongwlin/elf-aid-magic/internal/operator"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run elf-aid-magic.",
	Run:   runRun,
}

func runRun(_ *cobra.Command, _ []string) {
	stopped := false
	conf := config.New()

	l := logger.New(conf)
	defer l.Sync()

	l.Info("START")

	o := operator.New(conf, l)
	defer o.Destroy()

	fmt.Println("Link Start!")

	if !o.Connect() {
		fmt.Println("Failed to connect device.")
		o.Destroy()
		os.Exit(1)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		sig := <-sigs
		l.Info(
			"received interrupt signal to stop",
			zap.String("signal", sig.String()),
		)
		stopped = true
		cancel()
		o.Stop()
	}()

	go func() {
		defer wg.Done()
		if !o.Run(ctx) && !stopped {
			fmt.Println("Failed to run tasks.")
		}
	}()

	fmt.Println("Running...")

	wg.Wait()

	if stopped {
		fmt.Println("Interrupt")
	} else {
		fmt.Println("Completed")
	}
	l.Info("END")
}

func init() {
	rootCmd.AddCommand(runCmd)
}
