package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/service"
)

// GenericHandler provides HTTP handlers for generic service operations.
type GenericHandler[T any] struct {
	Service service.Interface[T]
}

// NewGenericHandler creates a new GenericHandler using the provided GenericRepository.
// It constructs the underlying GenericService from the repository.
func NewGenericHandler[T any](service service.Interface[T]) *GenericHandler[T] {
	return &GenericHandler[T]{
		Service: service,
	}
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
		if err := json.NewEncoder(w).Encode(created); err != nil {
			http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		}
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
		if err := json.NewEncoder(w).Encode(entity); err != nil {
			http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		}
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
		if err := json.NewEncoder(w).Encode(updated); err != nil {
			http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		}
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
		if err := json.NewEncoder(w).Encode(entities); err != nil {
			http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

// DeleteWhereHandler handles HTTP DELETE requests to delete entities matching a condition.
// It expects a JSON body with "condition" and "args".
func (h *GenericHandler[T]) DeleteWhereHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type DeleteWhereRequest struct {
			Condition string        `json:"condition"`
			Args      []interface{} `json:"args"`
		}
		var req DeleteWhereRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if err := h.Service.DeleteWhere(req.Condition, req.Args...); err != nil {
			http.Error(w, "Error deleting entities: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// FindHandler handles HTTP GET requests to find entities matching a condition.
// It expects a JSON body with "condition" and "args", and returns a JSON array of matching entities.
func (h *GenericHandler[T]) FindHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type FindRequest struct {
			Condition string        `json:"condition"`
			Args      []interface{} `json:"args"`
		}
		var req FindRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		entities, err := h.Service.Find(req.Condition, req.Args...)
		if err != nil {
			http.Error(w, "Error finding entities: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(entities); err != nil {
			http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

// FindFirstHandler handles HTTP GET requests to find the first entity matching a condition.
// It expects a JSON body with "condition" and "args", and returns the matching entity as JSON.
func (h *GenericHandler[T]) FindFirstHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type FindRequest struct {
			Condition string        `json:"condition"`
			Args      []interface{} `json:"args"`
		}
		var req FindRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		entity, err := h.Service.FindFirst(req.Condition, req.Args...)
		if err != nil {
			http.Error(w, "Error finding entity: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(entity); err != nil {
			http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

// CountHandler handles HTTP GET requests to count entities matching a condition.
// It expects a JSON body with "condition" and "args", and returns the count as JSON.
func (h *GenericHandler[T]) CountHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type CountRequest struct {
			Condition string        `json:"condition"`
			Args      []interface{} `json:"args"`
		}
		var req CountRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		count, err := h.Service.Count(req.Condition, req.Args...)
		if err != nil {
			http.Error(w, "Error counting entities: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]int64{"count": count}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

// GetPageHandler handles HTTP GET requests to retrieve a paginated list of entities.
// It expects a JSON body with "page", "pageSize", "condition", and "args", and returns a JSON array of entities.
func (h *GenericHandler[T]) GetPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type GetPageRequest struct {
			Page      int           `json:"page"`
			PageSize  int           `json:"pageSize"`
			Condition string        `json:"condition"`
			Args      []interface{} `json:"args"`
		}
		var req GetPageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		entities, err := h.Service.GetPage(req.Page, req.PageSize, req.Condition, req.Args...)
		if err != nil {
			http.Error(w, "Error retrieving page of entities: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(entities); err != nil {
			http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

// BulkInsertHandler handles HTTP POST requests to bulk insert multiple entities.
// It expects a JSON array of entities and returns a success message on completion.
func (h *GenericHandler[T]) BulkInsertHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var entities []T
		if err := json.NewDecoder(r.Body).Decode(&entities); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		// Convert []T to []*T
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

