package handler

import (
	"encoding/json"
	"net/http"

	"github.com/adal4ik/people-enrichment-service/internal/models"
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

func (p *PersonHandler) handleError(w http.ResponseWriter, r *http.Request, code int, message string, err error) {
	if err != nil {
		p.logger.Error(message,
			zap.Error(err),
			zap.Int("code", code),
			zap.String("url", r.URL.Path),
		)
	} else {
		p.logger.Error(message,
			zap.Int("code", code),
			zap.String("url", r.URL.Path),
		)
	}

	jsonErr := models.APIError{
		Code:     code,
		Message:  message,
		Resource: r.URL.Path,
	}
	jsonErr.Send(w)
}

func (p *PersonHandler) CreatePerson(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var person models.Person
	err := decoder.Decode(&person)
	if err != nil {
		p.handleError(w, req, 400, "failed to decode request body", err)
		return
	}
	p.logger.Debug("received person payload", zap.Any("person", person))

	err = p.service.CreatePerson(req.Context(), person)
	if err != nil {
		p.handleError(w, req, 500, "failed to save person", err)
		return
	}
	resp := models.APIResponse{
		Code:    201,
		Message: "Successfully created",
	}
	resp.Send(w)
}
