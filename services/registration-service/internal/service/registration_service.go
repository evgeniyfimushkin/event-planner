package service

import (
	"registration-service/internal/models"
	"registration-service/internal/repository"

	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/service"
)

// RegistrationService specializes in handling business logic for Registration entities.
// It embeds GenericService for basic CRUD operations and adds additional dependencies (e.g., a verifier).
type RegistrationService struct {
	*service.GenericService[models.Registration]
}

// NewRegistrationService creates a new instance of RegistrationService using the provided verifier and repository.
// It initializes the underlying GenericService using the given repository.
func NewRegistrationService(repo *repository.RegistrationRepository) *RegistrationService {
	return &RegistrationService{
		GenericService: service.NewGenericService[models.Registration](repo),
	}
}

