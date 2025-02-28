package repository

import (
	"fmt"
	"gorm.io/gorm"
)

type GenericRepository[T any] struct {
	db *gorm.DB
}

func NewGenericRepository[T any](db *gorm.DB) *GenericRepository[T] {
	return &GenericRepository[T]{db: db}
}

func (repo *GenericRepository[T]) Create(entity *T) error {
	result := repo.db.Create(entity)
	if result.Error != nil {
		return fmt.Errorf("unable to create entity: %w", result.Error)
	}
	return nil
}

func (repo *GenericRepository[T]) GetByID(id int) (*T, error) {
	var entity T
	result := repo.db.First(&entity, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to get entity by id: %w", result.Error)
	}
	return &entity, nil
}

func (repo *GenericRepository[T]) Update(entity *T) error {
	result := repo.db.Save(entity)
	if result.Error != nil {
		return fmt.Errorf("unable to update entity: %w", result.Error)
	}
	return nil
}

func (repo *GenericRepository[T]) Delete(id int) error {
	result := repo.db.Delete(new(T), id)
	if result.Error != nil {
		return fmt.Errorf("unable to delete entity: %w", result.Error)
	}
	return nil
}

func (repo *GenericRepository[T]) GetAll() ([]T, error) {
	var entities []T
	result := repo.db.Find(&entities)
	if result.Error != nil {
		return nil, fmt.Errorf("unable to fetch entities: %w", result.Error)
	}
	return entities, nil
}

