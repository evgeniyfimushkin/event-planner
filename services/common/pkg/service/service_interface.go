package service

// Interface defines all the generic service operations for type T.
type Interface[T any] interface {
	// Create creates a new entity.
	Create(entity *T) (*T, error)
	// GetByID retrieves an entity by its id.
	GetByID(id int) (*T, error)
	// Update updates an existing entity.
	Update(entity *T) (*T, error)
	// Delete deletes an entity by its id.
	Delete(id int) error
	// GetAll retrieves all entities.
	GetAll() ([]T, error)
	// DeleteWhere deletes entities matching the given condition.
	DeleteWhere(condition interface{}, args ...interface{}) error
	// Find returns all entities matching the given condition.
	Find(condition interface{}, args ...interface{}) ([]T, error)
	// FindFirst returns the first entity matching the given condition.
	FindFirst(condition interface{}, args ...interface{}) (*T, error)
	// Count returns the count of entities matching the given condition.
	Count(condition interface{}, args ...interface{}) (int64, error)
	// GetPage returns a paginated list of entities matching the given condition.
	GetPage(page int, pageSize int, condition interface{}, args ...interface{}) ([]T, error)
	// BulkInsert inserts multiple entities at once.
	BulkInsert(entities []*T) error
	// BulkUpdate updates multiple entities based on the given condition.
	BulkUpdate(condition interface{}, args []interface{}, updateData interface{}) error
}

