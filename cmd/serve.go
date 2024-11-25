package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dongwlin/elf-aid-magic/internal/config"
	"github.com/dongwlin/elf-aid-magic/internal/logger"
	"github.com/dongwlin/elf-aid-magic/internal/operator"
	"github.com/dongwlin/elf-aid-magic/internal/wire"
	"github.com/dongwlin/elf-aid-magic/public"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server at the specified address",
	Run:   serveRun,
}

func serveRun(cmd *cobra.Command, args []string) {
	conf := config.New()

	l := logger.New(conf)
	defer l.Sync()

	o := operator.New(conf, l)
	defer o.Destroy()

	h := wire.InitHandler(l, o)

	app := fiber.New()
	r := app.Group("/")
	initRouter(r, h)

	go func() {
		err := app.Listen(fmt.Sprintf(":%d", conf.Server.Port))
		if err != nil {
			o.Destroy()
			l.Error("failed to start server", zap.Error(err))
			fmt.Println("Failed to start server. See log.json for details.")
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	l.Info("shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		l.Fatal(
			"server shutdown",
			zap.Error(err),
		)
	}
	l.Info("server exiting")
}

func initRouter(r fiber.Router, h *wire.Handler) {
	r.Use(cors.New())

	r.Use(filesystem.New(filesystem.Config{
		Root:       http.FS(public.Public),
		PathPrefix: "dist",
	}))

	h.Ping.Register(r)

	ws := r.Group("/ws")
	ws.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	h.WebSocket.Register(ws)

	api := r.Group("/api")
	h.Vesrion.Register(api)
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
