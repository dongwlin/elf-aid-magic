package cmd

import (
	"fmt"
	"github.com/dongwlin/elf-aid-magic/internal/api"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server at the specified address",
	Run:   serveRun,
}

func serveRun(cmd *cobra.Command, args []string) {
	l := NewLogger()
	defer l.Sync()

	server := api.NewServer()
	server.SetupRouter()

	err := server.Start(":8000")
	if err != nil {
		l.Error("failed to start server", zap.Error(err))
		fmt.Println("Failed to start server. See log.json for details.")
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
