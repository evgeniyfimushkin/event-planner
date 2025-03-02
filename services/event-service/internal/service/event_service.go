package service

import (
	"event-service/internal/models"
	"event-service/internal/repository"
	"fmt"

	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/auth"
	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/logger"
)

type EventService struct {
    verifier *auth.Verifier
    eventRepo *repository.EventRepository
}

func NewEventService (verifier *auth.Verifier, eventRepo *repository.EventRepository) (*EventService, error) {
    if verifier == nil || eventRepo == nil {
        return nil, fmt.Errorf("failed to create EventService, arguments can not be nil")
    }
    return &EventService{
        verifier: verifier,
        eventRepo: eventRepo,
    }, nil
}

func (es *EventService) CreateEvent(accessToken string, event *models.Event) (*models.Event, error) {
}
