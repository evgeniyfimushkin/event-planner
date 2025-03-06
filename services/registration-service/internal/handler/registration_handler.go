package handler

import (
	"net/http"
	"registration-service/internal/models"
	"registration-service/internal/service"

	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/auth"
	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/handler"
)

type RegistrationHandler struct {
    *handler.GenericHandler[models.Registration]
}

func NewRegistrationHandler(service *service.RegistrationService, verifier *auth.Verifier) *RegistrationHandler {
    return &RegistrationHandler{
        GenericHandler: handler.NewGenericHandler[models.Registration](service, verifier),
    }
}

//func (h *RegistrationHandler) FindHandler() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		// Verify the JWT token and get claims.
//		claims, err := h.CheckToken(r)
//		if err != nil {
//			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
//			return
//		}
//
//		reserved := []string{"page", "pageSize"}
//		condition, args, err := parseQueryCondition(r.URL.Query(), reserved)
//		if err != nil {
//			http.Error(w, "Invalid query parameters: "+err.Error(), http.StatusBadRequest)
//			return
//		}
//
//		// Pass the claims to the service.
//		entities, err := h.Service.Find(claims, condition, args...)
//		if err != nil {
//			http.Error(w, "Error finding entities: "+err.Error(), http.StatusInternalServerError)
//			return
//		}
//
//		w.Header().Set("Content-Type", "application/json")
//		json.NewEncoder(w).Encode(entities)
//	}
//}
//
