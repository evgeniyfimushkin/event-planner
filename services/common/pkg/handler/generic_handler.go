package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/auth"
	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/service"
	"github.com/golang-jwt/jwt/v5"
)

// GenericHandler provides HTTP handlers for generic service operations.
type GenericHandler[T any] struct {
	Service  service.Interface[T]
	Verifier *auth.Verifier
}

// NewGenericHandler creates a new GenericHandler with the provided service and verifier.
func NewGenericHandler[T any](srv service.Interface[T], verif *auth.Verifier) *GenericHandler[T] {
	return &GenericHandler[T]{
		Service:  srv,
		Verifier: verif,
	}
}

// CheckToken extracts the JWT token from the "access_token" cookie and verifies it.
func (h *GenericHandler[T]) CheckToken(r *http.Request) (jwt.MapClaims, error) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		return nil, fmt.Errorf("missing access_token cookie: %w", err)
	}
	claims, err := h.Verifier.VerifyJWTToken(cookie.Value)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// ParseQueryCondition converts query parameters into an SQL-like condition and arguments.
func ParseQueryCondition(q url.Values, reserved []string) (string, []interface{}, error) {
	var conditions []string
	var args []interface{}

	// isReserved returns true if the key is in the reserved list.
	isReserved := func(key string) bool {
		for _, r := range reserved {
			if key == r {
				return true
			}
		}
		return false
	}

	for key, values := range q {
		if isReserved(key) {
			continue
		}
		if len(values) == 0 {
			continue
		}
		val := values[0]
		op := "="
		// Check for multi-character operators first.
		if len(val) >= 2 && (val[:2] == "<=" || val[:2] == ">=" || val[:2] == "!=") {
			op = val[:2]
			val = val[2:]
		} else if len(val) > 0 && (val[0] == '<' || val[0] == '>') {
			op = string(val[0])
			val = val[1:]
		}
		conditions = append(conditions, fmt.Sprintf("%s %s ?", key, op))
		args = append(args, val)
	}

	if len(conditions) == 0 {
		return "", nil, fmt.Errorf("no valid condition parameters provided")
	}

	condStr := conditions[0]
	for i := 1; i < len(conditions); i++ {
		condStr += " AND " + conditions[i]
	}

	return condStr, args, nil
}

// CreateHandler handles HTTP POST requests to create a new entity.
func (h *GenericHandler[T]) CreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify the JWT token and get claims.
		claims, err := h.CheckToken(r)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		var entity T
		if err := json.NewDecoder(r.Body).Decode(&entity); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Pass the claims to the service.
		created, err := h.Service.Create(claims, &entity)
		if err != nil {
			http.Error(w, "Error creating entity: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(created)
	}
}

// GetByIDHandler handles HTTP GET requests to retrieve an entity by its ID.
func (h *GenericHandler[T]) GetByIDHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify the JWT token and get claims.
		claims, err := h.CheckToken(r)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		idParam := r.URL.Query().Get("id")
		if idParam == "" {
			http.Error(w, "Missing id parameter", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Invalid id parameter", http.StatusBadRequest)
			return
		}

		// Pass the claims to the service.
		entity, err := h.Service.GetByID(claims, id)
		if err != nil {
			http.Error(w, "Entity not found: "+err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entity)
	}
}

// UpdateHandler handles HTTP PUT/PATCH requests to update an existing entity.
func (h *GenericHandler[T]) UpdateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify the JWT token and get claims.
		claims, err := h.CheckToken(r)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		var entity T
		if err := json.NewDecoder(r.Body).Decode(&entity); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Pass the claims to the service.
		updated, err := h.Service.Update(claims, &entity)
		if err != nil {
			http.Error(w, "Error updating entity: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updated)
	}
}

// DeleteHandler handles HTTP DELETE requests to delete an entity by its ID.
func (h *GenericHandler[T]) DeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify the JWT token and get claims.
		claims, err := h.CheckToken(r)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		idParam := r.URL.Query().Get("id")
		if idParam == "" {
			http.Error(w, "Missing id parameter", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Invalid id parameter", http.StatusBadRequest)
			return
		}

		// Pass the claims to the service.
		if err := h.Service.Delete(claims, id); err != nil {
			http.Error(w, "Error deleting entity: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// GetAllHandler handles HTTP GET requests to retrieve all entities.
func (h *GenericHandler[T]) GetAllHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify the JWT token and get claims.
		claims, err := h.CheckToken(r)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Pass the claims to the service.
		entities, err := h.Service.GetAll(claims)
		if err != nil {
			http.Error(w, "Error retrieving entities: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entities)
	}
}

// DeleteWhereHandler handles HTTP DELETE requests to delete entities matching a condition.
func (h *GenericHandler[T]) DeleteWhereHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify the JWT token and get claims.
		claims, err := h.CheckToken(r)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		reserved := []string{"page", "pageSize"}
		condition, args, err := ParseQueryCondition(r.URL.Query(), reserved)
		if err != nil {
			http.Error(w, "Invalid query parameters: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Pass the claims to the service.
		if err := h.Service.DeleteWhere(claims, condition, args...); err != nil {
			http.Error(w, "Error deleting entities: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// FindHandler handles HTTP GET requests to find entities matching a condition.
func (h *GenericHandler[T]) FindHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify the JWT token and get claims.
		claims, err := h.CheckToken(r)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		reserved := []string{"page", "pageSize"}
		condition, args, err := ParseQueryCondition(r.URL.Query(), reserved)
		if err != nil {
			http.Error(w, "Invalid query parameters: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Pass the claims to the service.
		entities, err := h.Service.Find(claims, condition, args...)
		if err != nil {
			http.Error(w, "Error finding entities: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entities)
	}
}

// FindFirstHandler handles HTTP GET requests to find the first entity matching a condition.
func (h *GenericHandler[T]) FindFirstHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify the JWT token and get claims.
		claims, err := h.CheckToken(r)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		reserved := []string{"page", "pageSize"}
		condition, args, err := ParseQueryCondition(r.URL.Query(), reserved)
		if err != nil {
			http.Error(w, "Invalid query parameters: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Pass the claims to the service.
		entity, err := h.Service.FindFirst(claims, condition, args...)
		if err != nil {
			http.Error(w, "Error finding entity: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entity)
	}
}

// CountHandler handles HTTP GET requests to count entities matching a condition.
func (h *GenericHandler[T]) CountHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify the JWT token and get claims.
		claims, err := h.CheckToken(r)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		reserved := []string{"page", "pageSize"}
		condition, args, err := ParseQueryCondition(r.URL.Query(), reserved)
		if err != nil {
			http.Error(w, "Invalid query parameters: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Pass the claims to the service.
		count, err := h.Service.Count(claims, condition, args...)
		if err != nil {
			http.Error(w, "Error counting entities: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int64{"count": count})
	}
}

// GetPageHandler handles HTTP GET requests to retrieve a paginated list of entities.
func (h *GenericHandler[T]) GetPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify the JWT token and get claims.
		claims, err := h.CheckToken(r)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		q := r.URL.Query()
		pageStr := q.Get("page")
		pageSizeStr := q.Get("pageSize")
		if pageStr == "" || pageSizeStr == "" {
			http.Error(w, "Missing page or pageSize parameters", http.StatusBadRequest)
			return
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, "Invalid page parameter", http.StatusBadRequest)
			return
		}

		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			http.Error(w, "Invalid pageSize parameter", http.StatusBadRequest)
			return
		}

		reserved := []string{"page", "pageSize"}
		condition, args, err := ParseQueryCondition(q, reserved)
		if err != nil {
			http.Error(w, "Invalid query parameters: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Pass the claims to the service.
		entities, err := h.Service.GetPage(claims, page, pageSize, condition, args...)
		if err != nil {
			http.Error(w, "Error retrieving page of entities: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entities)
	}
}

// BulkInsertHandler handles HTTP POST requests to bulk insert multiple entities.
func (h *GenericHandler[T]) BulkInsertHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify the JWT token and get claims.
		claims, err := h.CheckToken(r)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		var entities []T
		if err := json.NewDecoder(r.Body).Decode(&entities); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Convert []T to []*T.
		var entityPtrs []*T
		for i := range entities {
			entityPtrs = append(entityPtrs, &entities[i])
		}

		// Pass the claims to the service.
		if err := h.Service.BulkInsert(claims, entityPtrs); err != nil {
			http.Error(w, "Error bulk inserting entities: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Bulk insert successful"))
	}
}

// BulkUpdateHandler handles HTTP PUT requests to bulk update entities matching a condition.
func (h *GenericHandler[T]) BulkUpdateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify the JWT token and get claims.
		claims, err := h.CheckToken(r)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		type BulkUpdateRequest struct {
			Condition  string                 `json:"condition"`
			Args       []interface{}          `json:"args"`
			UpdateData map[string]interface{} `json:"updateData"`
		}
		var req BulkUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Pass the claims to the service.
		if err := h.Service.BulkUpdate(claims, req.Condition, req.Args, req.UpdateData); err != nil {
			http.Error(w, "Error bulk updating entities: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Bulk update successful"))
	}
}
