package service

import "github.com/golang-jwt/jwt/v5"

// Interface defines all the generic service operations for type T.
type Interface[T any] interface {
	// Create creates a new entity.
	Create(claims jwt.MapClaims, entity *T) (*T, error)
	// GetByID retrieves an entity by its id.
	GetByID(claims jwt.MapClaims, id int) (*T, error)
	// Update updates an existing entity.
	Update(claims jwt.MapClaims, entity *T) (*T, error)
	// Delete deletes an entity by its id.
	Delete(claims jwt.MapClaims, id int) error
	// GetAll retrieves all entities.
	GetAll(claims jwt.MapClaims, ) ([]T, error)
	// DeleteWhere deletes entities matching the given condition.
	DeleteWhere(claims jwt.MapClaims, condition interface{}, args ...interface{}) error
	// Find returns all entities matching the given condition.
	Find(claims jwt.MapClaims, condition interface{}, args ...interface{}) ([]T, error)
	// FindFirst returns the first entity matching the given condition.
	FindFirst(claims jwt.MapClaims, condition interface{}, args ...interface{}) (*T, error)
	// Count returns the count of entities matching the given condition.
	Count(claims jwt.MapClaims, condition interface{}, args ...interface{}) (int64, error)
	// GetPage returns a paginated list of entities matching the given condition.
	GetPage(claims jwt.MapClaims, page int, pageSize int, condition interface{}, args ...interface{}) ([]T, error)
	// BulkInsert inserts multiple entities at once.
	BulkInsert(claims jwt.MapClaims, entities []*T) error
	// BulkUpdate updates multiple entities based on the given condition.
	BulkUpdate(claims jwt.MapClaims, condition interface{}, args []interface{}, updateData interface{}) error
}

