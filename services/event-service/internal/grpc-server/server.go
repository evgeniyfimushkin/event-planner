package grpcserver

import (
	"context"
	"event-service/internal/service"

	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/protos/events"
	"google.golang.org/grpc"
)

type serverAPI struct {
    events.UnsafeEventServiceServer
    service *service.EventService
}

func Register(gRPC *grpc.Server, service *service.EventService) {
    events.RegisterEventServiceServer(gRPC, &serverAPI{service: service})
}



func (s *serverAPI) CheckAndReserve(ctx context.Context, req *events.CheckAndReserveRequest) (*events.CheckAndReserveResponse, error) {
    return &events.CheckAndReserveResponse{
        Status: events.ReserveStatus_SUCCESS,
    }, nil
}
