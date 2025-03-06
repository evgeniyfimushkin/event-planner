package service

import (
	"event-service/internal/models"
	"event-service/internal/repository"
	"fmt"
	"strings"
	"time"

	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/service"
	"github.com/golang-jwt/jwt/v5"
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

func (s *EventService) Create(claims jwt.MapClaims, entity *models.Event) (*models.Event, error) {
    username, ok := claims["username"].(string)
    if !ok {
        return nil, fmt.Errorf("invalid token: username not found or not a string")
    }
    entity.CreatedBy = username
    now := time.Now()
    oneYearLater := now.AddDate(1, 0, 0)

    if entity.StartTime.Before(now) {
        return nil, fmt.Errorf("start time cannot be in the past")
    }
    if entity.StartTime.After(oneYearLater) {
        return nil, fmt.Errorf("start time must be within one year from now")
    }
    if entity.EndTime.Before(entity.StartTime) {
        return nil, fmt.Errorf("end time must be after start time")
    }
    if entity.EndTime.After(oneYearLater) {
        return nil, fmt.Errorf("end time must be within one year from now")
    }

    if strings.TrimSpace(entity.Name) == "" {
        return nil, fmt.Errorf("name must not be empty")
    }

    if entity.MaxParticipants < 2 {
        return nil, fmt.Errorf("max participants must be at least 2")
    }

    entity.Participants = 1

    if entity.Participants > entity.MaxParticipants {
        return nil, fmt.Errorf("participants cannot exceed max participants")
    }

    if entity.Latitude < -90 || entity.Latitude > 90 {
        return nil, fmt.Errorf("latitude must be between -90 and 90")
    }
    if entity.Longitude < -180 || entity.Longitude > 180 {
        return nil, fmt.Errorf("longitude must be between -180 and 180")
    }

    return s.GenericService.Create(claims, entity)
}

func (s *EventService) Update(claims jwt.MapClaims, entity *models.Event) (*models.Event, error) {
    now := time.Now()
    oneYearLater := now.AddDate(1, 0, 0)

    if entity.StartTime.Before(now) {
        return nil, fmt.Errorf("start time cannot be in the past")
    }
    if entity.StartTime.After(oneYearLater) {
        return nil, fmt.Errorf("start time must be within one year from now")
    }
    if entity.EndTime.Before(entity.StartTime) {
        return nil, fmt.Errorf("end time must be after start time")
    }
    if entity.EndTime.After(oneYearLater) {
        return nil, fmt.Errorf("end time must be within one year from now")
    }

    if strings.TrimSpace(entity.Name) == "" {
        return nil, fmt.Errorf("name must not be empty")
    }

    if entity.MaxParticipants < 2 {
        return nil, fmt.Errorf("max participants must be at least 2")
    }

    if entity.Participants > entity.MaxParticipants {
        return nil, fmt.Errorf("participants cannot exceed max participants")
    }

    if entity.Latitude < -90 || entity.Latitude > 90 {
        return nil, fmt.Errorf("latitude must be between -90 and 90")
    }
    if entity.Longitude < -180 || entity.Longitude > 180 {
        return nil, fmt.Errorf("longitude must be between -180 and 180")
    }

    return s.GenericService.Update(claims, entity)
}
