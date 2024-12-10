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

var (
	id   string
	name string
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

	if id == "" {
		if len(conf.Taskers) == 0 {
			fmt.Println("taskers is empty")
			os.Exit(1)
		}
		id = conf.Taskers[0].ID
	}
	if name != "" {
		id = getTaskerIDByName(conf, name)
	}
	if id == "" {
		fmt.Println("tasker name not exists")
		os.Exit(1)
	}

	l.Info("START")

	fmt.Println("Link Start!")

	o := operator.New(conf, l, id)
	defer o.Destroy()

	initOperator(o)

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

func getTaskerIDByName(conf *config.Config, name string) string {
	for _, tasker := range conf.Taskers {
		if tasker.Name == name {
			return tasker.ID
		}
	}
	return ""
}

func initOperator(o *operator.Operator) {
	if !o.InitTasker() {
		fmt.Println("Failed to init takser.")
		o.Destroy()
		os.Exit(1)
	}

	if !o.InitResource() {
		fmt.Println("Failed to init resource.")
		o.Destroy()
		os.Exit(1)
	}

	if !o.InitController() {
		fmt.Println("Failed to init controller.")
		o.Destroy()
		os.Exit(1)
	}

	if !o.Connect() {
		fmt.Println("Failed to connect device.")
		o.Destroy()
		os.Exit(1)
	}
}

func init() {
	runCmd.PersistentFlags().StringVar(&id, "id", "", "Specify the tasker by id")
	runCmd.PersistentFlags().StringVar(&name, "name", "", "Specify the tasker by name")
	rootCmd.AddCommand(runCmd)
}
