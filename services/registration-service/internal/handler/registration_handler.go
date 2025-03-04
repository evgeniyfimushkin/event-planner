package handler

import (
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
