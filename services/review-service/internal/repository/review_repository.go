package repository

import (
	"review-service/internal/models"

	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/repository"
	"gorm.io/gorm"
)

type ReviewRepository struct {
    *repository.GenericRepository[models.Review]
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{
		GenericRepository: repository.NewGenericRepository[models.Review](db),
	}
}

// // GetByCategory returns all events that belong to the specified category.
// func (er *ReviewRepository) GetByCategory(category string) ([]models.Review, error) {
// 	var events []models.Review
// 	result := er.Db.Where("category = ?", category).Find(&events)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return events, nil
// }
// 
// // GetUpcomingReviews returns all events that have a start_time in the future.
// func (er *ReviewRepository) GetUpcomingReviews() ([]models.Review, error) {
// 	var events []models.Review
// 	result := er.Db.Where("start_time > NOW()").Order("start_time asc").Find(&events)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return events, nil
// }
