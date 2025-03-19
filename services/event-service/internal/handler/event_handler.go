package handler

import (
	"encoding/json"
	"event-service/internal/models"
	"event-service/internal/service"
	"net/http"
	"sort"

	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/auth"
	"github.com/evgeniyfimushkin/event-planner/services/common/pkg/handler"
)

type EventHandler struct {
    *handler.GenericHandler[models.Event]
    EventService *service.EventService
}

func NewEventHandler(service *service.EventService, verifier *auth.Verifier) *EventHandler {
    return &EventHandler{
        GenericHandler: handler.NewGenericHandler[models.Event](service, verifier),
        EventService: service,
    }
}

func (h *EventHandler) GetUpcomingEvents() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        claims, err := h.CheckToken(r)
        if err != nil {
            http.Error(w, "Unauthorized: " + err.Error(), http.StatusUnauthorized)
            return
        }

        events, err := h.EventService.GetUpcomingEvents(claims)
        if err != nil {
            http.Error(w, "Error retrieving events: " + err.Error(), http.StatusInternalServerError)
            return
        }
        sort.Slice(events, func(i, j int) bool {
           return events[i].StartTime.Before(events[j].StartTime)
        })
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(events)
    }
}

func (h *EventHandler) GetPreviousEvents() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        claims, err := h.CheckToken(r)
        if err != nil {
            http.Error(w, "Unauthorized: " + err.Error(), http.StatusUnauthorized)
            return
        }

        events, err := h.EventService.GetPreviousEvents(claims)
        if err != nil {
            http.Error(w, "Error retrieving events: " + err.Error(), http.StatusInternalServerError)
            return
        }

        sort.Slice(events, func(i, j int) bool {
            return events[i].StartTime.After(events[j].StartTime)
        })

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(events)
    }
}

