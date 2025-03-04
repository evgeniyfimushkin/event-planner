package repository

import (
	"registration-service/internal/models"

	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/repository"
	"gorm.io/gorm"
)

type RegistrationRepository struct {
    *repository.GenericRepository[models.Registration]
}

func NewRegistrationRepository(db *gorm.DB) *RegistrationRepository {
	return &RegistrationRepository{
		GenericRepository: repository.NewGenericRepository[models.Registration](db),
	}
}

// // GetByCategory returns all events that belong to the specified category.
// func (er *RegistrationRepository) GetByCategory(category string) ([]models.Registration, error) {
// 	var events []models.Registration
// 	result := er.Db.Where("category = ?", category).Find(&events)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return events, nil
// }
// 
// // GetUpcomingRegistrations returns all events that have a start_time in the future.
// func (er *RegistrationRepository) GetUpcomingRegistrations() ([]models.Registration, error) {
// 	var events []models.Registration
// 	result := er.Db.Where("start_time > NOW()").Order("start_time asc").Find(&events)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return events, nil
// }
