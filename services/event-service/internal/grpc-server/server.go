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
    event, err := s.service.GetByID(nil, int(req.EventId))
    if err != nil {
        return &events.CheckAndReserveResponse{
            Status: events.ReserveStatus_EVENT_NOT_FOUND,
        }, nil
    }

    if event.Participants >= event.MaxParticipants {
        return &events.CheckAndReserveResponse{
            Status: events.ReserveStatus_EVENT_FULL,
        }, nil
    }

    event.Participants = event.Participants + 1
    updatedEvent, err := s.service.Update(nil, event)
    if err != nil {
        return &events.CheckAndReserveResponse {
            Status: events.ReserveStatus_INTERNAL_ERROR,
        }, nil
    }

    return &events.CheckAndReserveResponse{
        Status: events.ReserveStatus_SUCCESS,
        CurrentParticipants: uint32(updatedEvent.Participants),
    }, nil
}
