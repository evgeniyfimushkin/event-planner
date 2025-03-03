package repository

import (
	"event-service/internal/models"

	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/repository"
	"gorm.io/gorm"
)

type EventRepository struct {
    *repository.GenericRepository[models.Event]
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{
		GenericRepository: repository.NewGenericRepository[models.Event](db),
	}
}

// // GetByCategory returns all events that belong to the specified category.
// func (er *EventRepository) GetByCategory(category string) ([]models.Event, error) {
// 	var events []models.Event
// 	result := er.Db.Where("category = ?", category).Find(&events)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return events, nil
// }
// 
// // GetUpcomingEvents returns all events that have a start_time in the future.
// func (er *EventRepository) GetUpcomingEvents() ([]models.Event, error) {
// 	var events []models.Event
// 	result := er.Db.Where("start_time > NOW()").Order("start_time asc").Find(&events)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return events, nil
// }
