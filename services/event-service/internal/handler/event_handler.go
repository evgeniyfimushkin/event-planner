package handler

import (
	"event-service/internal/models"
	"event-service/internal/service"

	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/handler"
)

type EventHandler struct {
    *handler.GenericHandler[models.Event]
}

func NewEventHandler(service *service.EventService) *EventHandler {
    return &EventHandler{
        GenericHandler: handler.NewGenericHandler[models.Event](service),
    }
}
