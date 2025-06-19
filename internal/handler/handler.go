package handler

import (
	"github.com/adal4ik/people-enrichment-service/internal/service"
	"go.uber.org/zap"
)

type Handler struct {
	PersonHandler *PersonHandler
}

func New(services *service.Service, logger *zap.Logger) *Handler {
	return &Handler{
		PersonHandler: NewPersonHandler(services.PersonService, logger),
	}
}
