package handler

import (
	"github.com/adal4ik/people-enrichment-service/internal/service"
	"go.uber.org/zap"
)

type PersonHandler struct {
	service service.PersonServiceInterface
	logger  *zap.Logger
}

func NewPersonHandler(service service.PersonServiceInterface, logger *zap.Logger) *PersonHandler {
	return &PersonHandler{
		service: service,
		logger:  logger,
	}
}
