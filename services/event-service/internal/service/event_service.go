package service

import (
	"event-service/internal/models"
	"event-service/internal/repository"
	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/service"
)

// EventService specializes in handling business logic for Event entities.
// It embeds GenericService for basic CRUD operations and adds additional dependencies (e.g., a verifier).
type EventService struct {
	*service.GenericService[models.Event]
}

// NewEventService creates a new instance of EventService using the provided verifier and repository.
// It initializes the underlying GenericService using the given repository.
func NewEventService(repo *repository.EventRepository) *EventService {
	return &EventService{
		GenericService: service.NewGenericService[models.Event](repo),
	}
}

