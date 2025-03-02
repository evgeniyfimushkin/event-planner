package repository

import (
	"fmt"
	"gorm.io/gorm"
)

type GenericRepository[T any] struct {
	Db *gorm.DB
}

func NewGenericRepository[T any](db *gorm.DB) *GenericRepository[T] {
	return &GenericRepository[T]{Db: db}
}

// Create creates a new entity and returns the created entity.
func (repo *GenericRepository[T]) Create(entity *T) (*T, error) {
	result := repo.Db.Create(entity)
	if result.Error != nil {
		return nil, fmt.Errorf("unable to create entity: %w", result.Error)
	}
	return entity, nil
}

// GetByID retrieves an entity by its id.
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

// Update updates an entity and returns the updated entity.
func (repo *GenericRepository[T]) Update(entity *T) (*T, error) {
	result := repo.Db.Save(entity)
	if result.Error != nil {
		return nil, fmt.Errorf("unable to update entity: %w", result.Error)
	}
	return entity, nil
}

// Delete deletes an entity by its id.
func (repo *GenericRepository[T]) Delete(id int) error {
	result := repo.Db.Delete(new(T), id)
	if result.Error != nil {
		return fmt.Errorf("unable to delete entity: %w", result.Error)
	}
	return nil
}

// GetAll retrieves all entities.
func (repo *GenericRepository[T]) GetAll() ([]T, error) {
	var entities []T
	result := repo.Db.Find(&entities)
	if result.Error != nil {
		return nil, fmt.Errorf("unable to fetch entities: %w", result.Error)
	}
	return entities, nil
}

// DeleteWhere deletes entities matching the given condition.
func (repo *GenericRepository[T]) DeleteWhere(condition interface{}, args ...interface{}) error {
	result := repo.Db.Where(condition, args...).Delete(new(T))
	if result.Error != nil {
		return fmt.Errorf("unable to delete entities with condition: %w", result.Error)
	}
	return nil
}

// Find returns all entities matching the given condition.
// Example: events, err := eventRepo.Find("category = ?", "Conference")
// Example: found, err := repo.Find("name LIKE ?", "Entity%")
func (repo *GenericRepository[T]) Find(condition interface{}, args ...interface{}) ([]T, error) {
	var entities []T
	result := repo.Db.Where(condition, args...).Find(&entities)
	if result.Error != nil {
		return nil, fmt.Errorf("unable to find entities: %w", result.Error)
	}
	return entities, nil
}

// FindFirst returns a single entity matching the given condition.
// Example: oneEntity, err := repo.FindFirst("name = ?", "Updated Entity")
func (repo *GenericRepository[T]) FindFirst(condition interface{}, args ...interface{}) (*T, error) {
	var entity T
	result := repo.Db.Where(condition, args...).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, result.Error
		}
		return nil, fmt.Errorf("unable to find entity: %w", result.Error)
	}
	return &entity, nil
}

// Count returns the count of entities matching the given condition.
// Example: count, err := eventRepo.Count("status = ? AND created_by = ?", "active", "john_doe")
// Example: count2, err := repo.Count("name = ? OR name = ?", "CountTest", "Other")
func (repo *GenericRepository[T]) Count(condition interface{}, args ...interface{}) (int64, error) {
	var count int64
	result := repo.Db.Model(new(T)).Where(condition, args...).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("unable to count entities: %w", result.Error)
	}
	return count, nil
}

// GetPage returns a paginated list of entities matching the given condition.
// Example: pageEvents, err := eventRepo.GetPage(1, 10, "city = ?", "San Francisco")
func (repo *GenericRepository[T]) GetPage(page int, pageSize int, condition interface{}, args ...interface{}) ([]T, error) {
	var entities []T
	offset := (page - 1) * pageSize
	result := repo.Db.Where(condition, args...).Offset(offset).Limit(pageSize).Find(&entities)
	if result.Error != nil {
		return nil, fmt.Errorf("unable to get page of entities: %w", result.Error)
	}
	return entities, nil
}

// BulkInsert inserts multiple entities at once.
// Example: eventRepo.BulkInsert(bulkEvents)
func (repo *GenericRepository[T]) BulkInsert(entities []*T) error {
	result := repo.Db.Create(&entities)
	if result.Error != nil {
		return fmt.Errorf("unable to bulk insert entities: %w", result.Error)
	}
	return nil
}

// BulkUpdate updates multiple entities based on the given condition with provided update data.
// err := eventRepo.BulkUpdate(
//     "category = ? AND created_by = ?",
//     []interface{}{"Conference", "john_doe"},
//     map[string]interface{}{"status": "inactive"},
// )
func (repo *GenericRepository[T]) BulkUpdate(condition interface{}, args []interface{}, updateData interface{}) error {
	result := repo.Db.Model(new(T)).Where(condition, args...).Updates(updateData)
	if result.Error != nil {
		return fmt.Errorf("unable to bulk update entities: %w", result.Error)
	}
	return nil
}


// ExecuteInTransaction executes the provided function within a transaction.
// If the function returns an error, the transaction is rolled back; otherwise, it is committed.
// Example:
// err := repo.ExecuteInTransaction(func(tx *gorm.DB) error {
// 		txRepo := NewGenericRepository[TestEntity](tx)
// 		entity, err := txRepo.Create(&TestEntity{Name: "Entity1"})
// 		if err != nil {
// 			return err
// 		}
// 		entity.Name = "Entity1 Updated"
// 		_, err = txRepo.Update(entity)
// 		if err != nil {
// 			return err
// 		}
//     }

func (repo *GenericRepository[T]) ExecuteInTransaction(fn func(tx *gorm.DB) error) error {
	tx := repo.Db.Begin()
	if tx.Error != nil {
		return tx.Error
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

