//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/dongwlin/elf-aid-magic/internal/handler"
	"github.com/dongwlin/elf-aid-magic/internal/logic"
	"github.com/dongwlin/elf-aid-magic/internal/operator"
	"github.com/google/wire"
	"go.uber.org/zap"
)

var logicSet = wire.NewSet(
	logic.NewPidLogic,
	logic.NewVersionLogic,
	logic.NewWebSocketLogic,
)

var handlerSet = wire.NewSet(
	handler.NewPidHandler,
	handler.NewPingHandler,
	handler.NewVersionHandler,
	handler.NewWebSocketHandler,
)

type Handler struct {
	Pid       *handler.PidHandler
	Ping      *handler.PingHandler
	Vesrion   *handler.VersionHandler
	WebSocket *handler.WebSocketHandler
}

func provideHandler(
	pidHandler *handler.PidHandler,
	pingHandler *handler.PingHandler,
	versionHandler *handler.VersionHandler,
	webSocketHandler *handler.WebSocketHandler,
) *Handler {
	return &Handler{
		Pid:       pidHandler,
		Ping:      pingHandler,
		Vesrion:   versionHandler,
		WebSocket: webSocketHandler,
	}
}

func InitHandler(logger *zap.Logger, om *operator.Manager) *Handler {
	wire.Build(logicSet, handlerSet, provideHandler)
	return nil
}
