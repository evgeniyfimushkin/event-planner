package service

import (
	"fmt"
	grpcclient "registration-service/internal/client/grpc-client"
	"registration-service/internal/models"
	"registration-service/internal/repository"
	"strconv"

	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/service"
	"github.com/golang-jwt/jwt/v5"
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
    if claims["userID"] != strconv.Itoa(int(entity.UserID)) {
        return entity, fmt.Errorf("UserID is not correct")
    }
    s.Find(

}
