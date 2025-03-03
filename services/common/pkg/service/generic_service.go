package service

import (
	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/repository"
	"gorm.io/gorm"
)

// GenericService provides a set of generic operations for entities of type T.
// It acts as a service layer that wraps around a repository implementation.
type GenericService[T any] struct {
	Repo repository.Interface[T]
}

// NewGenericService creates a new GenericService using the provided repository.
// It accepts any repository implementation that satisfies the repository.Interface[T] contract.
func NewGenericService[T any](repo repository.Interface[T]) *GenericService[T] {
	return &GenericService[T]{
		Repo: repo,
	}
}

// Create creates a new entity using the underlying repository.
// It returns the created entity or an error if the operation fails.
func (s *GenericService[T]) Create(entity *T) (*T, error) {
	created, err := s.Repo.Create(entity)
	if err != nil {
		return nil, err
	}
	return created, nil
}

// GetByID retrieves an entity by its unique identifier using the underlying repository.
// It returns the entity if found, or an error if not.
func (s *GenericService[T]) GetByID(id int) (*T, error) {
	return s.Repo.GetByID(id)
}

// Update updates an existing entity using the underlying repository.
// It returns the updated entity or an error if the update fails.
func (s *GenericService[T]) Update(entity *T) (*T, error) {
	updated, err := s.Repo.Update(entity)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

// Delete removes an entity identified by its unique identifier using the underlying repository.
// It returns an error if the deletion fails.
func (s *GenericService[T]) Delete(id int) error {
	return s.Repo.Delete(id)
}

// GetAll retrieves all entities from the underlying repository.
// It returns a slice of entities or an error if the retrieval fails.
func (s *GenericService[T]) GetAll() ([]T, error) {
	return s.Repo.GetAll()
}

// DeleteWhere deletes entities that match the specified condition using the underlying repository.
// The condition can be a SQL clause and args its corresponding parameters.
// It returns an error if the deletion fails.
func (s *GenericService[T]) DeleteWhere(condition interface{}, args ...interface{}) error {
	return s.Repo.DeleteWhere(condition, args...)
}

// Find returns all entities matching the specified condition using the underlying repository.
// It accepts a condition (e.g., SQL where clause) and variadic arguments for that condition,
// returning a slice of matching entities or an error if the query fails.
func (s *GenericService[T]) Find(condition interface{}, args ...interface{}) ([]T, error) {
	return s.Repo.Find(condition, args...)
}

// FindFirst returns the first entity that matches the specified condition using the underlying repository.
// It returns a pointer to the entity or an error if no matching entity is found.
func (s *GenericService[T]) FindFirst(condition interface{}, args ...interface{}) (*T, error) {
	return s.Repo.FindFirst(condition, args...)
}

// Count returns the number of entities that match the specified condition using the underlying repository.
// It returns the count as an int64 and any error encountered during the operation.
func (s *GenericService[T]) Count(condition interface{}, args ...interface{}) (int64, error) {
	return s.Repo.Count(condition, args...)
}

// GetPage retrieves a paginated list of entities that match the specified condition using the underlying repository.
// It accepts the page number, page size, a condition, and any arguments for that condition.
// It returns a slice of entities or an error if the operation fails.
func (s *GenericService[T]) GetPage(page int, pageSize int, condition interface{}, args ...interface{}) ([]T, error) {
	return s.Repo.GetPage(page, pageSize, condition, args...)
}

// BulkInsert inserts multiple entities at once using the underlying repository.
// It accepts a slice of pointers to entities and returns an error if the insertion fails.
func (s *GenericService[T]) BulkInsert(entities []*T) error {
	return s.Repo.BulkInsert(entities)
}

// BulkUpdate updates multiple entities that match the specified condition using the underlying repository.
// It accepts a condition, arguments for that condition, and a map containing the update data.
// It returns an error if the update operation fails.
func (s *GenericService[T]) BulkUpdate(condition interface{}, args []interface{}, updateData interface{}) error {
	return s.Repo.BulkUpdate(condition, args, updateData)
}

// ExecuteInTransaction executes the provided function within a database transaction using the underlying repository.
// If the function returns an error, the transaction is rolled back; otherwise, it is committed.
// It returns any error encountered during the transactional operation.
func (s *GenericService[T]) ExecuteInTransaction(fn func(tx *gorm.DB) error) error {
	return s.Repo.ExecuteInTransaction(fn)
}

