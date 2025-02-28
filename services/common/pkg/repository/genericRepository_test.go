package repository

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestEntity struct {
	ID   int
	Name string
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err, "failed to open in-memory database")

	err = db.AutoMigrate(&TestEntity{})
	assert.NoError(t, err, "failed to migrate test entity")

	return db
}

func TestGenericRepository_CreateAndGetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGenericRepository[TestEntity](db)

	entity := &TestEntity{Name: "TestName"}
	err := repo.Create(entity)
	assert.NoError(t, err, "failed to create entity")
	assert.NotZero(t, entity.ID, "entity ID should be set after creation")

	result, err := repo.GetByID(entity.ID)
	assert.NoError(t, err, "failed to get entity by ID")
	assert.NotNil(t, result, "result should not be nil")
	assert.Equal(t, entity.Name, result.Name, "retrieved entity name should match")
}

func TestGenericRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGenericRepository[TestEntity](db)

	entity := &TestEntity{Name: "OldName"}
	err := repo.Create(entity)
	assert.NoError(t, err)

	entity.Name = "NewName"
	err = repo.Update(entity)
	assert.NoError(t, err)

	result, err := repo.GetByID(entity.ID)
	assert.NoError(t, err)
	assert.Equal(t, "NewName", result.Name, "entity name should be updated")
}

func TestGenericRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGenericRepository[TestEntity](db)

	entity := &TestEntity{Name: "ToBeDeleted"}
	err := repo.Create(entity)
	assert.NoError(t, err)

	err = repo.Delete(entity.ID)
	assert.NoError(t, err)

	result, err := repo.GetByID(entity.ID)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	assert.Nil(t, result, "entity should be deleted")
}

func TestGenericRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGenericRepository[TestEntity](db)

	repo.Create(&TestEntity{Name: "A"})
	repo.Create(&TestEntity{Name: "B"})

	entities, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, entities, 2, "should retrieve all entities")
}

