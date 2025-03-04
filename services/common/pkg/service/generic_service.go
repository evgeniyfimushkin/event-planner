package service

import (
	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/repository"
	"github.com/golang-jwt/jwt/v5"
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
// It ignores the claims parameter.
func (s *GenericService[T]) Create(claims jwt.MapClaims, entity *T) (*T, error) {
	return s.Repo.Create(entity)
}

// GetByID retrieves an entity by its unique identifier using the underlying repository.
// It ignores the claims parameter.
func (s *GenericService[T]) GetByID(claims jwt.MapClaims, id int) (*T, error) {
	return s.Repo.GetByID(id)
}

// Update updates an existing entity using the underlying repository.
// It ignores the claims parameter.
func (s *GenericService[T]) Update(claims jwt.MapClaims, entity *T) (*T, error) {
	return s.Repo.Update(entity)
}

// Delete removes an entity identified by its unique identifier using the underlying repository.
// It ignores the claims parameter.
func (s *GenericService[T]) Delete(claims jwt.MapClaims, id int) error {
	return s.Repo.Delete(id)
}

// GetAll retrieves all entities from the underlying repository.
// It ignores the claims parameter.
func (s *GenericService[T]) GetAll(claims jwt.MapClaims) ([]T, error) {
	return s.Repo.GetAll()
}

// DeleteWhere deletes entities that match the specified condition using the underlying repository.
// It ignores the claims parameter.
func (s *GenericService[T]) DeleteWhere(claims jwt.MapClaims, condition interface{}, args ...interface{}) error {
	return s.Repo.DeleteWhere(condition, args...)
}

// Find returns all entities matching the specified condition using the underlying repository.
// It ignores the claims parameter.
func (s *GenericService[T]) Find(claims jwt.MapClaims, condition interface{}, args ...interface{}) ([]T, error) {
	return s.Repo.Find(condition, args...)
}

// FindFirst returns the first entity matching the specified condition using the underlying repository.
// It ignores the claims parameter.
func (s *GenericService[T]) FindFirst(claims jwt.MapClaims, condition interface{}, args ...interface{}) (*T, error) {
	return s.Repo.FindFirst(condition, args...)
}

// Count returns the number of entities that match the specified condition using the underlying repository.
// It ignores the claims parameter.
func (s *GenericService[T]) Count(claims jwt.MapClaims, condition interface{}, args ...interface{}) (int64, error) {
	return s.Repo.Count(condition, args...)
}

// GetPage retrieves a paginated list of entities that match the specified condition using the underlying repository.
// It ignores the claims parameter.
func (s *GenericService[T]) GetPage(claims jwt.MapClaims, page int, pageSize int, condition interface{}, args ...interface{}) ([]T, error) {
	return s.Repo.GetPage(page, pageSize, condition, args...)
}

// BulkInsert inserts multiple entities at once using the underlying repository.
// It ignores the claims parameter.
func (s *GenericService[T]) BulkInsert(claims jwt.MapClaims, entities []*T) error {
	return s.Repo.BulkInsert(entities)
}

// BulkUpdate updates multiple entities that match the specified condition using the underlying repository.
// It ignores the claims parameter.
func (s *GenericService[T]) BulkUpdate(claims jwt.MapClaims, condition interface{}, args []interface{}, updateData interface{}) error {
	return s.Repo.BulkUpdate(condition, args, updateData)
}
