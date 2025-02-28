package repository

import (
	"fmt"
	"gorm.io/gorm"
)

// GenericRepository is generic wrapper for gorm.DB interface
// that has standard CRUD functions
type GenericRepository[T any] struct {
	Db *gorm.DB
}

// NewGenericRepository return pointer to new GenericRepository
func NewGenericRepository[T any](db *gorm.DB) *GenericRepository[T] {
	return &GenericRepository[T]{Db: db}
}

// Create creates new entity
func (repo *GenericRepository[T]) Create(entity *T) error {
	result := repo.Db.Create(entity)
	if result.Error != nil {
		return fmt.Errorf("unable to create entity: %w", result.Error)
	}
	return nil
}

// GetByID gets pointer to entity by id int 
func (repo *GenericRepository[T]) GetByID(id int) (*T, error) {
	var entity T
	result := repo.Db.First(&entity, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, result.Error
		}
		return nil, fmt.Errorf("unable to get entity by id: %w", result.Error)
	}
	return &entity, nil
}

// Update updates one entity
func (repo *GenericRepository[T]) Update(entity *T) error {
	result := repo.Db.Save(entity)
	if result.Error != nil {
		return fmt.Errorf("unable to update entity: %w", result.Error)
	}
	return nil
}

// Delete deletes entity by id
func (repo *GenericRepository[T]) Delete(id int) error {
	result := repo.Db.Delete(new(T), id)
	if result.Error != nil {
		return fmt.Errorf("unable to delete entity: %w", result.Error)
	}
	return nil
}

// GetAll get all entities from this repository
func (repo *GenericRepository[T]) GetAll() ([]T, error) {
	var entities []T
	result := repo.Db.Find(&entities)
	if result.Error != nil {
		return nil, fmt.Errorf("unable to fetch entities: %w", result.Error)
	}
	return entities, nil
}

