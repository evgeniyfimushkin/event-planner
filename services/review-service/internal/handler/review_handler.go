package handler

import (
	"encoding/json"
	"net/http"
	"review-service/internal/models"
	"review-service/internal/service"

	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/auth"
	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/handler"
)

type ReviewHandler struct {
    *handler.GenericHandler[models.Review]
}

func NewReviewHandler(service *service.ReviewService, verifier *auth.Verifier) *ReviewHandler {
    return &ReviewHandler{
        GenericHandler: handler.NewGenericHandler[models.Review](service, verifier),
    }
}

func (h *ReviewHandler) GetMyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify the JWT token and get claims.
		claims, err := h.CheckToken(r)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

        userIDFloat, ok := claims["userID"].(float64)
        if !ok {
            http.Error(w, "UserID is not a number", http.StatusUnauthorized)
            return
        }
    
        userID := uint(userIDFloat) 

        condition := "user_id = ?"
        args := []interface{}{userID}

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

func (h *ReviewHandler) DeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := h.CheckToken(r)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		var req struct {
			EventID uint `json:"event_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = h.Service.Delete(claims, int(req.EventID))
		if err != nil {
			http.Error(w, "Failed to delete Review: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

