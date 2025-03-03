package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/handler"
	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/repository"
	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestEntity is a dummy struct for testing GenericHandler.
type TestEntity struct {
	ID   int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name"`
}

// setupTestDB creates an isolated in-memory SQLite database and migrates TestEntity.
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// Limit the connection pool to force a single connection (ensuring isolation for :memory: DB).
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get generic database: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)

	if err := db.AutoMigrate(&TestEntity{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

// setupHandler creates a new GenericHandler[TestEntity] instance using a GenericRepository.
func setupHandler(t *testing.T) (*handler.GenericHandler[TestEntity], *gorm.DB) {
	db := setupTestDB(t)
	repo := repository.NewGenericRepository[TestEntity](db)
    service := service.NewGenericService(repo)
	h := handler.NewGenericHandler(service)
	return h, db
}

// TestCreateHandler tests the CreateHandler endpoint.
func TestCreateHandler(t *testing.T) {
	h, _ := setupHandler(t)
	entity := TestEntity{Name: "TestCreate"}
	body, _ := json.Marshal(entity)
	req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateHandler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var created TestEntity
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}
	if created.ID == 0 {
		t.Fatalf("expected valid ID, got %d", created.ID)
	}
	if created.Name != entity.Name {
		t.Fatalf("expected name %s, got %s", entity.Name, created.Name)
	}
}

// TestGetByIDHandler tests the GetByIDHandler endpoint.
func TestGetByIDHandler(t *testing.T) {
	h, db := setupHandler(t)
	// Create entity directly in DB.
	entity := TestEntity{Name: "TestGetByID"}
	if err := db.Create(&entity).Error; err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/get?id="+strconv.Itoa(entity.ID), nil)
	w := httptest.NewRecorder()
	h.GetByIDHandler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var fetched TestEntity
	if err := json.Unmarshal(w.Body.Bytes(), &fetched); err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}
	if fetched.ID != entity.ID {
		t.Fatalf("expected ID %d, got %d", entity.ID, fetched.ID)
	}
}

// TestUpdateHandler tests the UpdateHandler endpoint.
func TestUpdateHandler(t *testing.T) {
	h, db := setupHandler(t)
	// Create entity.
	entity := TestEntity{Name: "Original"}
	if err := db.Create(&entity).Error; err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}
	// Update entity.
	entity.Name = "Updated"
	body, _ := json.Marshal(entity)
	req := httptest.NewRequest(http.MethodPut, "/update", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.UpdateHandler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var updated TestEntity
	if err := json.Unmarshal(w.Body.Bytes(), &updated); err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}
	if updated.Name != "Updated" {
		t.Fatalf("expected name Updated, got %s", updated.Name)
	}
}

// TestDeleteHandler tests the DeleteHandler endpoint.
func TestDeleteHandler(t *testing.T) {
	h, db := setupHandler(t)
	// Create entity.
	entity := TestEntity{Name: "ToDelete"}
	if err := db.Create(&entity).Error; err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}
	req := httptest.NewRequest(http.MethodDelete, "/delete?id="+strconv.Itoa(entity.ID), nil)
	w := httptest.NewRecorder()
	h.DeleteHandler().ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", w.Code)
	}
	// Verify that entity is deleted.
	var count int64
	db.Model(&TestEntity{}).Where("id = ?", entity.ID).Count(&count)
	if count != 0 {
		t.Fatalf("expected entity to be deleted, count %d", count)
	}
}

// TestGetAllHandler tests the GetAllHandler endpoint.
func TestGetAllHandler(t *testing.T) {
	h, db := setupHandler(t)
	// Create multiple entities.
	entities := []TestEntity{
		{Name: "Entity1"},
		{Name: "Entity2"},
	}
	for _, e := range entities {
		if err := db.Create(&e).Error; err != nil {
			t.Fatalf("failed to create entity: %v", err)
		}
	}
	req := httptest.NewRequest(http.MethodGet, "/all", nil)
	w := httptest.NewRecorder()
	h.GetAllHandler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var fetched []TestEntity
	if err := json.Unmarshal(w.Body.Bytes(), &fetched); err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}
	if len(fetched) != 2 {
		t.Fatalf("expected 2 entities, got %d", len(fetched))
	}
}

// TestDeleteWhereHandler tests the DeleteWhereHandler endpoint.
func TestDeleteWhereHandler(t *testing.T) {
	h, db := setupHandler(t)
	// Create entities: 2 with name "DeleteTest" and 1 with "KeepTest".
	entities := []TestEntity{
		{Name: "DeleteTest"},
		{Name: "DeleteTest"},
		{Name: "KeepTest"},
	}
	for _, e := range entities {
		if err := db.Create(&e).Error; err != nil {
			t.Fatalf("failed to create entity: %v", err)
		}
	}
	// Delete where name = "DeleteTest".
	reqBody := map[string]interface{}{
		"condition": "name = ?",
		"args":      []interface{}{"DeleteTest"},
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodDelete, "/deletewhere", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.DeleteWhereHandler().ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", w.Code)
	}
	// Verify that only the entity with "KeepTest" remains.
	var remaining []TestEntity
	if err := db.Find(&remaining).Error; err != nil {
		t.Fatalf("failed to fetch entities: %v", err)
	}
	if len(remaining) != 1 || remaining[0].Name != "KeepTest" {
		t.Fatalf("unexpected remaining entities: %+v", remaining)
	}
}

// TestFindHandler tests the FindHandler endpoint.
func TestFindHandler(t *testing.T) {
	h, db := setupHandler(t)
	// Create entities.
	entities := []TestEntity{
		{Name: "FindTest"},
		{Name: "FindTest"},
		{Name: "OtherTest"},
	}
	for _, e := range entities {
		if err := db.Create(&e).Error; err != nil {
			t.Fatalf("failed to create entity: %v", err)
		}
	}
	reqBody := map[string]interface{}{
		"condition": "name = ?",
		"args":      []interface{}{"FindTest"},
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodGet, "/find", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.FindHandler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var found []TestEntity
	if err := json.Unmarshal(w.Body.Bytes(), &found); err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}
	if len(found) != 2 {
		t.Fatalf("expected 2 entities, got %d", len(found))
	}
}

// TestFindFirstHandler tests the FindFirstHandler endpoint.
func TestFindFirstHandler(t *testing.T) {
	h, db := setupHandler(t)
	// Create several entities with the same name.
	entities := []TestEntity{
		{Name: "FirstTest"},
		{Name: "FirstTest"},
	}
	for _, e := range entities {
		if err := db.Create(&e).Error; err != nil {
			t.Fatalf("failed to create entity: %v", err)
		}
	}
	reqBody := map[string]interface{}{
		"condition": "name = ?",
		"args":      []interface{}{"FirstTest"},
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodGet, "/findfirst", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.FindFirstHandler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var found TestEntity
	if err := json.Unmarshal(w.Body.Bytes(), &found); err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}
	if found.Name != "FirstTest" {
		t.Fatalf("expected name FirstTest, got %s", found.Name)
	}
}

// TestCountHandler tests the CountHandler endpoint.
func TestCountHandler(t *testing.T) {
	h, db := setupHandler(t)
	// Create entities.
	entities := []TestEntity{
		{Name: "CountTest"},
		{Name: "CountTest"},
		{Name: "OtherTest"},
	}
	for _, e := range entities {
		if err := db.Create(&e).Error; err != nil {
			t.Fatalf("failed to create entity: %v", err)
		}
	}
	reqBody := map[string]interface{}{
		"condition": "name = ?",
		"args":      []interface{}{"CountTest"},
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodGet, "/count", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CountHandler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var resp map[string]int64
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}
	if count, ok := resp["count"]; !ok || count != 2 {
		t.Fatalf("expected count 2, got %v", resp)
	}
}

// TestGetPageHandler tests the GetPageHandler endpoint.
func TestGetPageHandler(t *testing.T) {
	h, db := setupHandler(t)
	// Create 5 entities.
	for i := 1; i <= 5; i++ {
		e := TestEntity{Name: "PageTest" + strconv.Itoa(i)}
		if err := db.Create(&e).Error; err != nil {
			t.Fatalf("failed to create entity: %v", err)
		}
	}
	reqBody := map[string]interface{}{
		"page":      1,
		"pageSize":  2,
		"condition": "name LIKE ?",
		"args":      []interface{}{"PageTest%"},
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodGet, "/getpage", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.GetPageHandler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var page []TestEntity
	if err := json.Unmarshal(w.Body.Bytes(), &page); err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}
	if len(page) != 2 {
		t.Fatalf("expected 2 entities, got %d", len(page))
	}
}

// TestBulkInsertHandler tests the BulkInsertHandler endpoint.
func TestBulkInsertHandler(t *testing.T) {
	h, db := setupHandler(t)
	entities := []TestEntity{
		{Name: "Bulk1"},
		{Name: "Bulk2"},
	}
	body, _ := json.Marshal(entities)
	req := httptest.NewRequest(http.MethodPost, "/bulkinsert", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.BulkInsertHandler().ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}
	var count int64
	db.Model(&TestEntity{}).Count(&count)
	if count != 2 {
		t.Fatalf("expected 2 entities, got %d", count)
	}
}

// TestBulkUpdateHandler tests the BulkUpdateHandler endpoint.
func TestBulkUpdateHandler(t *testing.T) {
	h, db := setupHandler(t)
	// Create entities.
	entities := []TestEntity{
		{Name: "BulkUpdateTest"},
		{Name: "BulkUpdateTest"},
		{Name: "Other"},
	}
	for _, e := range entities {
		if err := db.Create(&e).Error; err != nil {
			t.Fatalf("failed to create entity: %v", err)
		}
	}
	reqBody := map[string]interface{}{
		"condition":  "name = ?",
		"args":       []interface{}{"BulkUpdateTest"},
		"updateData": map[string]interface{}{"name": "UpdatedBulk"},
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/bulkupdate", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.BulkUpdateHandler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var count int64
	db.Model(&TestEntity{}).Where("name = ?", "UpdatedBulk").Count(&count)
	if count != 2 {
		t.Fatalf("expected 2 updated entities, got %d", count)
	}
}

