package repository

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestEntity — test entity for demonstration of work genericRepository
type TestEntity struct {
	ID   int `gorm:"primaryKey"`
	Name string
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err, "failed to open in-memory database")
	err = db.AutoMigrate(&TestEntity{})
	assert.NoError(t, err, "failed to migrate TestEntity")
	return db
}

func TestGenericRepository_CreateAndGetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGenericRepository[TestEntity](db)

	entity := &TestEntity{Name: "TestEntity"}
	created, err := repo.Create(entity)
	assert.NoError(t, err, "failed to create entity")
	assert.NotZero(t, created.ID, "entity ID should be set after creation")

	fetched, err := repo.GetByID(created.ID)
	assert.NoError(t, err, "failed to get entity by ID")
	assert.NotNil(t, fetched, "result should not be nil")
	assert.Equal(t, created.Name, fetched.Name, "retrieved entity name should match")
}

func TestGenericRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGenericRepository[TestEntity](db)

	entity := &TestEntity{Name: "OldName"}
	created, err := repo.Create(entity)
	assert.NoError(t, err)

	created.Name = "NewName"
	updated, err := repo.Update(created)
	assert.NoError(t, err, "failed to update entity")
	assert.Equal(t, "NewName", updated.Name, "entity name should be updated")

	fetched, err := repo.GetByID(created.ID)
	assert.NoError(t, err)
	assert.Equal(t, "NewName", fetched.Name, "updated entity name should match")
}

func TestGenericRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGenericRepository[TestEntity](db)

	entity := &TestEntity{Name: "DeleteTest"}
	created, err := repo.Create(entity)
	assert.NoError(t, err)

	err = repo.Delete(created.ID)
	assert.NoError(t, err, "failed to delete entity")

	fetched, err := repo.GetByID(created.ID)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), "expected ErrRecordNotFound")
	assert.Nil(t, fetched, "entity should be deleted")
}

func TestGenericRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGenericRepository[TestEntity](db)

	_, err := repo.Create(&TestEntity{Name: "A"})
	assert.NoError(t, err)
	_, err = repo.Create(&TestEntity{Name: "B"})
	assert.NoError(t, err)

	all, err := repo.GetAll()
	assert.NoError(t, err, "failed to fetch all entities")
	assert.Len(t, all, 2, "should retrieve all entities")
}

func TestGenericRepository_DeleteWhere(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGenericRepository[TestEntity](db)

	_, err := repo.Create(&TestEntity{Name: "ToDelete"})
	assert.NoError(t, err)
	_, err = repo.Create(&TestEntity{Name: "ToDelete"})
	assert.NoError(t, err)
	_, err = repo.Create(&TestEntity{Name: "NotToDelete"})
	assert.NoError(t, err)

	err = repo.DeleteWhere("name = ?", "ToDelete")
	assert.NoError(t, err, "failed to delete entities with condition")

	remaining, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, remaining, 1, "only one entity should remain")
	assert.Equal(t, "NotToDelete", remaining[0].Name, "remaining entity should be 'NotToDelete'")
}

func TestGenericRepository_Find(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGenericRepository[TestEntity](db)

	_, err := repo.Create(&TestEntity{Name: "Entity1"})
	assert.NoError(t, err)
	_, err = repo.Create(&TestEntity{Name: "Entity2"})
	assert.NoError(t, err)
	_, err = repo.Create(&TestEntity{Name: "Other"})
	assert.NoError(t, err)

	found, err := repo.Find("name LIKE ?", "Entity%")
	assert.NoError(t, err, "failed to find entities")
	assert.Len(t, found, 2, "should find 2 entities")
}

func TestGenericRepository_FindFirst(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGenericRepository[TestEntity](db)

	_, err := repo.Create(&TestEntity{Name: "UniqueEntity"})
	assert.NoError(t, err)

	entity, err := repo.FindFirst("name = ?", "UniqueEntity")
	assert.NoError(t, err, "failed to find one entity")
	assert.NotNil(t, entity)
	assert.Equal(t, "UniqueEntity", entity.Name, "entity name should match")

	// Проверяем случай, когда сущность не найдена.
	entity, err = repo.FindFirst("name = ?", "NonExistent")
	assert.Error(t, err, "expected error for non-existent entity")
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), "error should be ErrRecordNotFound")
	assert.Nil(t, entity)
}

func TestGenericRepository_Count(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGenericRepository[TestEntity](db)

	_, err := repo.Create(&TestEntity{Name: "CountTest"})
	assert.NoError(t, err)
	_, err = repo.Create(&TestEntity{Name: "CountTest"})
	assert.NoError(t, err)
	_, err = repo.Create(&TestEntity{Name: "Other"})
	assert.NoError(t, err)

	count, err := repo.Count("name = ?", "CountTest")
	assert.NoError(t, err, "failed to count entities")
	assert.Equal(t, int64(2), count, "expected count to be 2")

	// Пример с несколькими плейсхолдерами.
	count2, err := repo.Count("name = ? OR name = ?", "CountTest", "Other")
	assert.NoError(t, err, "failed to count entities with OR condition")
	assert.Equal(t, int64(3), count2, "expected count to be 3")
}

func TestGenericRepository_GetPage(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGenericRepository[TestEntity](db)

	for i := 1; i <= 15; i++ {
		_, err := repo.Create(&TestEntity{Name: fmt.Sprintf("Entity %d", i)})
		assert.NoError(t, err)
	}

	page1, err := repo.GetPage(1, 10, "name LIKE ?", "Entity%")
	assert.NoError(t, err, "failed to get page 1")
	assert.Len(t, page1, 10, "page 1 should contain 10 entities")

	page2, err := repo.GetPage(2, 10, "name LIKE ?", "Entity%")
	assert.NoError(t, err, "failed to get page 2")
	assert.Len(t, page2, 5, "page 2 should contain 5 entities")
}

func TestGenericRepository_BulkInsert(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGenericRepository[TestEntity](db)

	entities := []*TestEntity{
		{Name: "Bulk1"},
		{Name: "Bulk2"},
		{Name: "Bulk3"},
	}

	err := repo.BulkInsert(entities)
	assert.NoError(t, err, "failed to bulk insert entities")

	all, err := repo.GetAll()
	assert.NoError(t, err, "failed to get all entities after bulk insert")
	assert.Len(t, all, 3, "should have 3 entities after bulk insert")
}

func TestGenericRepository_BulkUpdate(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGenericRepository[TestEntity](db)

	_, err := repo.Create(&TestEntity{Name: "OldName1"})
	assert.NoError(t, err)
	_, err = repo.Create(&TestEntity{Name: "OldName2"})
	assert.NoError(t, err)
	_, err = repo.Create(&TestEntity{Name: "Other"})
	assert.NoError(t, err)

	// Обновляем все записи, у которых имя начинается с "OldName"
	err = repo.BulkUpdate("name LIKE ?", []interface{}{"OldName%"}, map[string]interface{}{"name": "NewName"})
	assert.NoError(t, err, "failed to bulk update entities")

	updated, err := repo.Find("name = ?", "NewName")
	assert.NoError(t, err, "failed to find updated entities")
	assert.Len(t, updated, 2, "expected 2 entities to be updated")
}

