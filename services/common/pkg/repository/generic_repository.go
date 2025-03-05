package repository

import (
	"fmt"
	"gorm.io/gorm"
)

// GenericRepository - generic struct, that provides all necessary methods to work with GORM
type GenericRepository[T any] struct {
	Db *gorm.DB
}

// NewGenericRepository - constructor for NewGenericRepository
func NewGenericRepository[T any](db *gorm.DB) *GenericRepository[T] {
	return &GenericRepository[T]{Db: db}
}

// Create creates a new entity and returns the created entity.
func (repo *GenericRepository[T]) Create(entity *T) (*T, error) {
	result := repo.Db.Create(entity)
	if result.Error != nil {
		return nil, fmt.Errorf("%w: %v", ErrCreateEntity, result.Error)
	}
	return entity, nil
}

// GetByID retrieves an entity by its id.
func (repo *GenericRepository[T]) GetByID(id int) (*T, error) {
	var entity T
	result := repo.Db.First(&entity, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrGetEntityByID, result.Error)
	}
	return &entity, nil
}

// Update updates an entity and returns the updated entity.
func (repo *GenericRepository[T]) Update(entity *T) (*T, error) {
	result := repo.Db.Save(entity)
	if result.Error != nil {
		return nil, fmt.Errorf("%w: %v", ErrUpdateEntity, result.Error)
	}
	return entity, nil
}

// Delete deletes an entity by its id.
func (repo *GenericRepository[T]) Delete(id int) error {
	result := repo.Db.Delete(new(T), id)
	if result.Error != nil {
		return fmt.Errorf("%w: %v", ErrDeleteEntity, result.Error)
	}
	return nil
}

// GetAll retrieves all entities.
func (repo *GenericRepository[T]) GetAll() ([]T, error) {
	var entities []T
	result := repo.Db.Find(&entities)
	if result.Error != nil {
		return nil, fmt.Errorf("%w: %v", ErrFetchEntities, result.Error)
	}
	return entities, nil
}

// DeleteWhere deletes entities matching the given condition.
func (repo *GenericRepository[T]) DeleteWhere(condition interface{}, args ...interface{}) error {
	result := repo.Db.Where(condition, args...).Delete(new(T))
	if result.Error != nil {
		return fmt.Errorf("%w: %v", ErrDeleteWithCond, result.Error)
	}
	return nil
}

// Find returns all entities matching the given condition.
func (repo *GenericRepository[T]) Find(condition interface{}, args ...interface{}) ([]T, error) {
	var entities []T
	result := repo.Db.Where(condition, args...).Find(&entities)
	if result.Error != nil {
		return nil, fmt.Errorf("%w: %v", ErrFindEntities, result.Error)
	}
	return entities, nil
}

// FindFirst returns a single entity matching the given condition.
func (repo *GenericRepository[T]) FindFirst(condition interface{}, args ...interface{}) (*T, error) {
	var entity T
	result := repo.Db.Where(condition, args...).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrFindEntity, result.Error)
	}
	return &entity, nil
}

// Count returns the count of entities matching the given condition.
func (repo *GenericRepository[T]) Count(condition interface{}, args ...interface{}) (int64, error) {
	var count int64
	result := repo.Db.Model(new(T)).Where(condition, args...).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("%w: %v", ErrCountEntities, result.Error)
	}
	return count, nil
}

// GetPage returns a paginated list of entities matching the given condition.
func (repo *GenericRepository[T]) GetPage(page int, pageSize int, condition interface{}, args ...interface{}) ([]T, error) {
	var entities []T
	offset := (page - 1) * pageSize
	result := repo.Db.Where(condition, args...).Offset(offset).Limit(pageSize).Find(&entities)
	if result.Error != nil {
		return nil, fmt.Errorf("%w: %v", ErrGetPageEntities, result.Error)
	}
	return entities, nil
}

// BulkInsert inserts multiple entities at once.
func (repo *GenericRepository[T]) BulkInsert(entities []*T) error {
	result := repo.Db.Create(&entities)
	if result.Error != nil {
		return fmt.Errorf("%w: %v", ErrBulkInsert, result.Error)
	}
	return nil
}

// BulkUpdate updates multiple entities based on the given condition with provided update data.
func (repo *GenericRepository[T]) BulkUpdate(condition interface{}, args []interface{}, updateData interface{}) error {
	result := repo.Db.Model(new(T)).Where(condition, args...).Updates(updateData)
	if result.Error != nil {
		return fmt.Errorf("%w: %v", ErrBulkUpdate, result.Error)
	}
	return nil
}

// ExecuteInTransaction executes the provided function within a transaction.
func (repo *GenericRepository[T]) ExecuteInTransaction(fn func(tx *gorm.DB) error) error {
	tx := repo.Db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("%w: %v", ErrTransaction, tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

