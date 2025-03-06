package service

import (
	"context"
	"fmt"
	grpcclient "registration-service/internal/client/grpc-client"
	"registration-service/internal/models"
	"registration-service/internal/repository"

	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/protos/events"
	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/service"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// RegistrationService specializes in handling business logic for Registration entities.
// It embeds GenericService for basic CRUD operations and adds additional dependencies (e.g., a verifier).
type RegistrationService struct {
	*service.GenericService[models.Registration]
    eventClient *grpcclient.EventClient
}

// NewRegistrationService creates a new instance of RegistrationService using the provided verifier and repository.
// It initializes the underlying GenericService using the given repository.
func NewRegistrationService(repo *repository.RegistrationRepository, grpcclient *grpcclient.EventClient) *RegistrationService {
	return &RegistrationService{
		GenericService: service.NewGenericService[models.Registration](repo),
        eventClient: grpcclient,
	}
}

func (s *RegistrationService) Create(claims jwt.MapClaims, entity *models.Registration) (*models.Registration, error) {
    userIDFloat, ok := claims["userID"].(float64)
    if !ok {
        return nil, fmt.Errorf("UserID is not a number")
    }

    entity.UserID = uint(userIDFloat) 

    existing, err := s.FindFirst(claims, "event_id = ? AND user_id = ?", entity.EventID, entity.UserID)
    if err == nil && existing != nil {
        return nil, fmt.Errorf("user is already registered for this event")
    }
    
    username, ok := claims["username"].(string)
    if !ok {
        return nil, fmt.Errorf("invalid token: username not found or not a string")
    }

    if err != nil && err != gorm.ErrRecordNotFound {
        return nil, fmt.Errorf("error checking existing registration: %w", err)
    }
    
    resp, err := s.eventClient.CheckAndReserve(context.Background(), uint32(entity.EventID), username)
    if err != nil {
        return nil, err
    }
    if resp.Status == events.ReserveStatus_RESERVE_STATUS_UNSPECIFIED {
        return nil, fmt.Errorf("Event creator cannot register for their own event")
    }

    if resp.Status == events.ReserveStatus_EVENT_NOT_FOUND {
        return nil, fmt.Errorf("Event with id %d not found",entity.EventID)
    }

    if resp.Status == events.ReserveStatus_EVENT_FULL {
        return nil, fmt.Errorf("Event with id %d is full", entity.EventID)
    }

    if resp.Status == events.ReserveStatus_INTERNAL_ERROR {
        return nil, fmt.Errorf("Internal error")
    }

    updatedRegistration , err := s.GenericService.Create(claims, entity)
    if err != nil {
        return nil, err
    }

    return updatedRegistration, nil
}

func (s *RegistrationService) Delete(claims jwt.MapClaims, id int) error {

    userIDFloat, ok := claims["userID"].(float64)
    if !ok {
        return fmt.Errorf("UserID is not a number")
    }

    userID := uint(userIDFloat) 
    username, ok := claims["username"].(string)
    if !ok {
        return fmt.Errorf("invalid token: username not found or not a string")
    }
 

    existing, err := s.FindFirst(claims, "event_id = ? AND user_id = ?", id, userID)
    if err != nil || existing == nil {
        return fmt.Errorf("You're not registrated")
    }

    if existing.UserID != userID {
        return fmt.Errorf("It's not your registration!")
    }

       
    resp, err := s.eventClient.RemoveRegistration(context.Background(), uint32(id), username)
    if err != nil {
        return err
    }
    if resp.Status == events.ReserveStatus_EVENT_NOT_FOUND {
        return fmt.Errorf("Event with id %d not found", id)
    }

    if resp.Status == events.ReserveStatus_INTERNAL_ERROR {
        return fmt.Errorf("Internal error")
    }

    err = s.GenericService.Delete(claims, int(existing.ID))
    if err != nil {
        return err
    }

    return nil
}

