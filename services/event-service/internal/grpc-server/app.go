package grpcserver

import (
	"event-service/internal/service"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type App struct {
    log         *slog.Logger
    gRPCServcer *grpc.Server
    port        int
}

func New(
    service *service.EventService,
    log *slog.Logger,
    port int,
) *App {
    gRPCServer := grpc.NewServer()
    Register(gRPCServer, service)
    return &App {
        log: log,
        gRPCServcer: gRPCServer,
        port: port,
    }
}

func (a *App) MustRun() {
    if err := a.Run(); err != nil {
        panic(err)
    }
}

func (a *App) Run() error {
    a.log.Info("starting gRPC server")
    l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
    if err != nil {
        return err
    }
    a.log.Info(fmt.Sprintf("grpc server is running on port %d", a.port))
    if err := a.gRPCServcer.Serve(l); err != nil {
        return err
    }
    return nil
}

func (a *App) Stop() {
    a.log.Info("stopping gRPC server")
    a.gRPCServcer.GracefulStop()
}
