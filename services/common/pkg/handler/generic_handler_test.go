package handler_test

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/auth"
	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/handler"
	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/repository"
	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/service"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestEntity is the test entity used for integration tests.
type TestEntity struct {
	ID   int    `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

// generateTestKey creates a new ECDSA key pair and returns the private key,
// public key, and a Base64-encoded PEM string of the public key.
func generateTestKey(t *testing.T) (*ecdsa.PrivateKey, *ecdsa.PublicKey, string) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	pub := &priv.PublicKey
	pubBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		t.Fatalf("failed to marshal public key: %v", err)
	}
	pemBlock := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})
	pubKeyStr := base64.StdEncoding.EncodeToString(pemBlock)
	return priv, pub, pubKeyStr
}

// generateValidToken creates a signed JWT token with a one-hour expiration using the given private key.
func generateValidToken(t *testing.T, priv *ecdsa.PrivateKey) string {
	claims := jwt.MapClaims{
		"user": "test",
		"exp":  time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenStr, err := token.SignedString(priv)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	return tokenStr
}

// setupTestDB opens an in-memory SQLite database and migrates the TestEntity model.
func setupTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	if err := db.AutoMigrate(&TestEntity{}); err != nil {
		t.Fatalf("failed to migrate TestEntity: %v", err)
	}
	return db
}

// newTestHandler creates a repository, service, and handler for TestEntity.
// It returns the handler, the test DB, and the private key used to sign tokens.
func newTestHandler(t *testing.T) (*handler.GenericHandler[TestEntity], *gorm.DB, *ecdsa.PrivateKey) {
	db := setupTestDB(t)
	repo := repository.NewGenericRepository[TestEntity](db)
	svc := service.NewGenericService[TestEntity](repo)
	// Generate a real key pair.
	priv, _, pubKeyStr := generateTestKey(t)
	verif, err := auth.NewVerifier(pubKeyStr)
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}
	h := handler.NewGenericHandler[TestEntity](svc, verif)
	return h, db, priv
}

// addValidCookie adds a cookie named "access_token" with the given token to the request.
func addValidCookie(r *http.Request, token string) {
	r.AddCookie(&http.Cookie{
		Name:  "access_token",
		Value: token,
	})
}

// ---------------------- Integration Tests ----------------------

// TestCreateHandlerIntegration tests the CreateHandler.
func TestCreateHandlerIntegration(t *testing.T) {
	h, db, priv := newTestHandler(t)

	reqBody := `{"name": "test entity"}`
	req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewBufferString(reqBody))
	token := generateValidToken(t, priv)
	addValidCookie(req, token)
	rec := httptest.NewRecorder()

	h.CreateHandler()(rec, req)
	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	var entity TestEntity
	if err := json.NewDecoder(res.Body).Decode(&entity); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if entity.ID == 0 {
		t.Errorf("expected non-zero ID, got %d", entity.ID)
	}
	if entity.Name != "test entity" {
		t.Errorf("expected name 'test entity', got '%s'", entity.Name)
	}

	// Verify record in database.
	var dbEntity TestEntity
	if err := db.First(&dbEntity, entity.ID).Error; err != nil {
		t.Errorf("failed to find entity in DB: %v", err)
	}
}

// TestGetByIDHandlerIntegration tests the GetByIDHandler.
func TestGetByIDHandlerIntegration(t *testing.T) {
	h, db, priv := newTestHandler(t)

	// Pre-create an entity.
	entity := TestEntity{Name: "getByID test"}
	if err := db.Create(&entity).Error; err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}

	urlStr := fmt.Sprintf("/get?id=%d", entity.ID)
	req := httptest.NewRequest(http.MethodGet, urlStr, nil)
	token := generateValidToken(t, priv)
	addValidCookie(req, token)
	rec := httptest.NewRecorder()

	h.GetByIDHandler()(rec, req)
	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	var result TestEntity
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result.ID != entity.ID || result.Name != entity.Name {
		t.Errorf("unexpected result: %+v", result)
	}
}

// TestUpdateHandlerIntegration tests the UpdateHandler.
func TestUpdateHandlerIntegration(t *testing.T) {
	h, db, priv := newTestHandler(t)

	// Create an entity.
	entity := TestEntity{Name: "original"}
	if err := db.Create(&entity).Error; err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}

	updatedBody := fmt.Sprintf(`{"id": %d, "name": "updated"}`, entity.ID)
	req := httptest.NewRequest(http.MethodPut, "/update", bytes.NewBufferString(updatedBody))
	token := generateValidToken(t, priv)
	addValidCookie(req, token)
	rec := httptest.NewRecorder()

	h.UpdateHandler()(rec, req)
	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}
	var updatedEntity TestEntity
	if err := json.NewDecoder(res.Body).Decode(&updatedEntity); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if updatedEntity.Name != "updated" {
		t.Errorf("expected name 'updated', got '%s'", updatedEntity.Name)
	}

	// Verify update in database.
	var dbEntity TestEntity
	if err := db.First(&dbEntity, entity.ID).Error; err != nil {
		t.Fatalf("failed to find updated entity in DB: %v", err)
	}
	if dbEntity.Name != "updated" {
		t.Errorf("DB not updated, got '%s'", dbEntity.Name)
	}
}

// TestDeleteHandlerIntegration tests the DeleteHandler.
func TestDeleteHandlerIntegration(t *testing.T) {
	h, db, priv := newTestHandler(t)

	// Create an entity.
	entity := TestEntity{Name: "to delete"}
	if err := db.Create(&entity).Error; err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}

	urlStr := fmt.Sprintf("/delete?id=%d", entity.ID)
	req := httptest.NewRequest(http.MethodDelete, urlStr, nil)
	token := generateValidToken(t, priv)
	addValidCookie(req, token)
	rec := httptest.NewRecorder()

	h.DeleteHandler()(rec, req)
	res := rec.Result()
	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", res.StatusCode)
	}

	// Verify deletion in database.
	var dbEntity TestEntity
	err := db.First(&dbEntity, entity.ID).Error
	if err == nil {
		t.Errorf("entity still exists in DB after deletion")
	}
}

// TestGetAllHandlerIntegration tests the GetAllHandler.
func TestGetAllHandlerIntegration(t *testing.T) {
	h, db, priv := newTestHandler(t)

	// Create several entities.
	entities := []TestEntity{
		{Name: "A"},
		{Name: "B"},
	}
	for _, e := range entities {
		if err := db.Create(&e).Error; err != nil {
			t.Fatalf("failed to create entity: %v", err)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/all", nil)
	token := generateValidToken(t, priv)
	addValidCookie(req, token)
	rec := httptest.NewRecorder()

	h.GetAllHandler()(rec, req)
	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	var results []TestEntity
	if err := json.NewDecoder(res.Body).Decode(&results); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(results) < 2 {
		t.Errorf("expected at least 2 entities, got %d", len(results))
	}
}

// TestDeleteWhereHandlerIntegration tests the DeleteWhereHandler.
func TestDeleteWhereHandlerIntegration(t *testing.T) {
	h, db, priv := newTestHandler(t)

	// Create entities with different Name values.
	entities := []TestEntity{
		{Name: "delete"},
		{Name: "keep"},
		{Name: "delete"},
	}
	for _, e := range entities {
		if err := db.Create(&e).Error; err != nil {
			t.Fatalf("failed to create entity: %v", err)
		}
	}

	// Delete records where name equals "delete".
	q := url.Values{}
	q.Set("name", "delete")
	req := httptest.NewRequest(http.MethodDelete, "/deleteWhere?"+q.Encode(), nil)
	token := generateValidToken(t, priv)
	addValidCookie(req, token)
	rec := httptest.NewRecorder()

	h.DeleteWhereHandler()(rec, req)
	res := rec.Result()
	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", res.StatusCode)
	}

	// Verify remaining records.
	var remaining []TestEntity
	if err := db.Find(&remaining).Error; err != nil {
		t.Fatalf("failed to query DB: %v", err)
	}
	for _, r := range remaining {
		if r.Name == "delete" {
			t.Errorf("entity with name 'delete' still exists")
		}
	}
}

// TestFindHandlerIntegration tests the FindHandler.
func TestFindHandlerIntegration(t *testing.T) {
	h, db, priv := newTestHandler(t)

	// Create entities.
	entities := []TestEntity{
		{Name: "findme"},
		{Name: "other"},
		{Name: "findme"},
	}
	for _, e := range entities {
		if err := db.Create(&e).Error; err != nil {
			t.Fatalf("failed to create entity: %v", err)
		}
	}

	q := url.Values{}
	q.Set("name", "findme")
	req := httptest.NewRequest(http.MethodGet, "/find?"+q.Encode(), nil)
	token := generateValidToken(t, priv)
	addValidCookie(req, token)
	rec := httptest.NewRecorder()

	h.FindHandler()(rec, req)
	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}
	var results []TestEntity
	if err := json.NewDecoder(res.Body).Decode(&results); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

// TestFindFirstHandlerIntegration tests the FindFirstHandler.
func TestFindFirstHandlerIntegration(t *testing.T) {
	h, db, priv := newTestHandler(t)

	// Create two entities with the same name.
	entity := TestEntity{Name: "first"}
	if err := db.Create(&entity).Error; err != nil {
		t.Fatalf("failed to create entity: %v", err)
	}
	_ = db.Create(&TestEntity{Name: "first"})

	q := url.Values{}
	q.Set("name", "first")
	req := httptest.NewRequest(http.MethodGet, "/findFirst?"+q.Encode(), nil)
	token := generateValidToken(t, priv)
	addValidCookie(req, token)
	rec := httptest.NewRecorder()

	h.FindFirstHandler()(rec, req)
	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}
	var result TestEntity
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result.Name != "first" {
		t.Errorf("expected name 'first', got '%s'", result.Name)
	}
}

// TestCountHandlerIntegration tests the CountHandler.
func TestCountHandlerIntegration(t *testing.T) {
	h, db, priv := newTestHandler(t)

	// Create 5 entities with the same name.
	for i := 0; i < 5; i++ {
		if err := db.Create(&TestEntity{Name: "count"}).Error; err != nil {
			t.Fatalf("failed to create entity: %v", err)
		}
	}

	q := url.Values{}
	q.Set("name", "count")
	req := httptest.NewRequest(http.MethodGet, "/count?"+q.Encode(), nil)
	token := generateValidToken(t, priv)
	addValidCookie(req, token)
	rec := httptest.NewRecorder()

	h.CountHandler()(rec, req)
	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}
	var result map[string]int64
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	count, ok := result["count"]
	if !ok || count != 5 {
		t.Errorf("expected count 5, got %v", result)
	}
}

// TestGetPageHandlerIntegration tests the GetPageHandler.
func TestGetPageHandlerIntegration(t *testing.T) {
	h, db, priv := newTestHandler(t)

	// Create 15 entities with the same name.
	for i := 0; i < 15; i++ {
		if err := db.Create(&TestEntity{Name: "page"}).Error; err != nil {
			t.Fatalf("failed to create entity: %v", err)
		}
	}

	q := url.Values{}
	q.Set("page", "2")
	q.Set("pageSize", "5")
	q.Set("name", "page")
	req := httptest.NewRequest(http.MethodGet, "/page?"+q.Encode(), nil)
	token := generateValidToken(t, priv)
	addValidCookie(req, token)
	rec := httptest.NewRecorder()

	h.GetPageHandler()(rec, req)
	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}
	var results []TestEntity
	if err := json.NewDecoder(res.Body).Decode(&results); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(results) != 5 {
		t.Errorf("expected 5 entities, got %d", len(results))
	}
}

// TestBulkInsertHandlerIntegration tests the BulkInsertHandler.
func TestBulkInsertHandlerIntegration(t *testing.T) {
	h, db, priv := newTestHandler(t)

	reqBody := `[{"name": "bulk1"}, {"name": "bulk2"}]`
	req := httptest.NewRequest(http.MethodPost, "/bulkInsert", bytes.NewBufferString(reqBody))
	token := generateValidToken(t, priv)
	addValidCookie(req, token)
	rec := httptest.NewRecorder()

	h.BulkInsertHandler()(rec, req)
	res := rec.Result()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", res.StatusCode)
	}
	body, _ := io.ReadAll(res.Body)
	if string(body) != "Bulk insert successful" {
		t.Errorf("unexpected response body: %s", string(body))
	}
	var count int64
	if err := db.Model(&TestEntity{}).Count(&count).Error; err != nil {
		t.Fatalf("failed to count entities: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 entities, got %d", count)
	}
}

// TestBulkUpdateHandlerIntegration tests the BulkUpdateHandler.
func TestBulkUpdateHandlerIntegration(t *testing.T) {
	h, db, priv := newTestHandler(t)

	// Create 3 entities with the name "old".
	for i := 0; i < 3; i++ {
		if err := db.Create(&TestEntity{Name: "old"}).Error; err != nil {
			t.Fatalf("failed to create entity: %v", err)
		}
	}

	reqBody := `{"condition": "name = ?", "args": ["old"], "updateData": {"name": "new"}}`
	req := httptest.NewRequest(http.MethodPut, "/bulkUpdate", bytes.NewBufferString(reqBody))
	token := generateValidToken(t, priv)
	addValidCookie(req, token)
	rec := httptest.NewRecorder()

	h.BulkUpdateHandler()(rec, req)
	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}
	body, _ := io.ReadAll(res.Body)
	if string(body) != "Bulk update successful" {
		t.Errorf("unexpected response body: %s", string(body))
	}
	var entities []TestEntity
	if err := db.Where("name = ?", "new").Find(&entities).Error; err != nil {
		t.Fatalf("failed to query updated entities: %v", err)
	}
	if len(entities) != 3 {
		t.Errorf("expected 3 updated entities, got %d", len(entities))
	}
}

// TestMissingTokenIntegration tests that a request without the access_token cookie is rejected.
func TestMissingTokenIntegration(t *testing.T) {
	h, _, _ := newTestHandler(t)

	reqBody := `{"name": "no token"}`
	req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewBufferString(reqBody))
	// Do not add the access_token cookie.
	rec := httptest.NewRecorder()

	h.CreateHandler()(rec, req)
	res := rec.Result()
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401 for missing token, got %d", res.StatusCode)
	}
}

