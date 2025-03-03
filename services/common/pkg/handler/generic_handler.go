package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/service"
)

// GenericHandler provides HTTP handlers for generic service operations.
type GenericHandler[T any] struct {
	Service service.Interface[T]
}

// NewGenericHandler creates a new GenericHandler using the provided service.
func NewGenericHandler[T any](srv service.Interface[T]) *GenericHandler[T] {
	return &GenericHandler[T]{
		Service: srv,
	}
}

// parseQueryCondition parses query parameters into a SQL-like condition string and corresponding arguments.
// It ignores reserved keys (such as "id", "page", "pageSize").
// For each query parameter, if the value begins with an operator (<=, >=, !=, <, >), that operator is used;
// otherwise, "=" is assumed.
func parseQueryCondition(q url.Values, reserved []string) (string, []interface{}, error) {
	var conditions []string
	var args []interface{}

	// isReserved returns true if key is in the reserved list.
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

	// Join conditions with " AND ".
	condStr := conditions[0]
	for i := 1; i < len(conditions); i++ {
		condStr += " AND " + conditions[i]
	}

	return condStr, args, nil
}

// CreateHandler handles HTTP POST requests to create a new entity.
// It decodes the request body into an entity and returns the created entity as JSON.
func (h *GenericHandler[T]) CreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var entity T
		if err := json.NewDecoder(r.Body).Decode(&entity); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		created, err := h.Service.Create(&entity)
		if err != nil {
			http.Error(w, "Error creating entity: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(created)
	}
}

// GetByIDHandler handles HTTP GET requests to retrieve an entity by its ID.
// The entity ID is expected as a query parameter "id".
func (h *GenericHandler[T]) GetByIDHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		entity, err := h.Service.GetByID(id)
		if err != nil {
			http.Error(w, "Entity not found: "+err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entity)
	}
}

// UpdateHandler handles HTTP PUT/PATCH requests to update an existing entity.
// It decodes the request body into an entity and returns the updated entity as JSON.
func (h *GenericHandler[T]) UpdateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var entity T
		if err := json.NewDecoder(r.Body).Decode(&entity); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		updated, err := h.Service.Update(&entity)
		if err != nil {
			http.Error(w, "Error updating entity: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updated)
	}
}

// DeleteHandler handles HTTP DELETE requests to delete an entity by its ID.
// The entity ID is expected as a query parameter "id". On success, it returns a No Content status.
func (h *GenericHandler[T]) DeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		if err := h.Service.Delete(id); err != nil {
			http.Error(w, "Error deleting entity: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// GetAllHandler handles HTTP GET requests to retrieve all entities.
// It returns a JSON array of all entities.
func (h *GenericHandler[T]) GetAllHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entities, err := h.Service.GetAll()
		if err != nil {
			http.Error(w, "Error retrieving entities: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entities)
	}
}

// DeleteWhereHandler handles HTTP DELETE requests to delete entities matching a condition.
// The condition is constructed from query parameters (excluding reserved keys such as "id", "page", "pageSize").
// If query-to-condition parsing fails, an error is returned.
func (h *GenericHandler[T]) DeleteWhereHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reserved := []string{"id", "page", "pageSize"}
		condition, args, err := parseQueryCondition(r.URL.Query(), reserved)
		if err != nil {
			http.Error(w, "Invalid query parameters: "+err.Error(), http.StatusBadRequest)
			return
		}
		if err := h.Service.DeleteWhere(condition, args...); err != nil {
			http.Error(w, "Error deleting entities: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// FindHandler handles HTTP GET requests to find entities matching a condition.
// The condition is constructed from query parameters (excluding reserved keys), and returns a JSON array of matching entities.
func (h *GenericHandler[T]) FindHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reserved := []string{"id", "page", "pageSize"}
		condition, args, err := parseQueryCondition(r.URL.Query(), reserved)
		if err != nil {
			http.Error(w, "Invalid query parameters: "+err.Error(), http.StatusBadRequest)
			return
		}
		entities, err := h.Service.Find(condition, args...)
		if err != nil {
			http.Error(w, "Error finding entities: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entities)
	}
}

// FindFirstHandler handles HTTP GET requests to find the first entity matching a condition.
// The condition is constructed from query parameters (excluding reserved keys), and returns the matching entity as JSON.
func (h *GenericHandler[T]) FindFirstHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reserved := []string{"id", "page", "pageSize"}
		condition, args, err := parseQueryCondition(r.URL.Query(), reserved)
		if err != nil {
			http.Error(w, "Invalid query parameters: "+err.Error(), http.StatusBadRequest)
			return
		}
		entity, err := h.Service.FindFirst(condition, args...)
		if err != nil {
			http.Error(w, "Error finding entity: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entity)
	}
}

// CountHandler handles HTTP GET requests to count entities matching a condition.
// The condition is constructed from query parameters (excluding reserved keys), and returns the count as JSON.
func (h *GenericHandler[T]) CountHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reserved := []string{"id", "page", "pageSize"}
		condition, args, err := parseQueryCondition(r.URL.Query(), reserved)
		if err != nil {
			http.Error(w, "Invalid query parameters: "+err.Error(), http.StatusBadRequest)
			return
		}
		count, err := h.Service.Count(condition, args...)
		if err != nil {
			http.Error(w, "Error counting entities: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int64{"count": count})
	}
}

// GetPageHandler handles HTTP GET requests to retrieve a paginated list of entities matching a condition.
// It expects "page" and "pageSize" query parameters along with additional parameters for conditions.
func (h *GenericHandler[T]) GetPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		reserved := []string{"id", "page", "pageSize"}
		condition, args, err := parseQueryCondition(q, reserved)
		if err != nil {
			http.Error(w, "Invalid query parameters: "+err.Error(), http.StatusBadRequest)
			return
		}

		entities, err := h.Service.GetPage(page, pageSize, condition, args...)
		if err != nil {
			http.Error(w, "Error retrieving page of entities: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entities)
	}
}

// BulkInsertHandler handles HTTP POST requests to bulk insert multiple entities.
// It expects a JSON array of entities in the request body and returns a success message on completion.
func (h *GenericHandler[T]) BulkInsertHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		if err := h.Service.BulkInsert(entityPtrs); err != nil {
			http.Error(w, "Error bulk inserting entities: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Bulk insert successful"))
	}
}

// BulkUpdateHandler handles HTTP PUT requests to bulk update entities matching a condition.
// It expects a JSON body with "condition", "args", and "updateData", and returns a success message.
func (h *GenericHandler[T]) BulkUpdateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		if err := h.Service.BulkUpdate(req.Condition, req.Args, req.UpdateData); err != nil {
			http.Error(w, "Error bulk updating entities: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Bulk update successful"))
	}
}

