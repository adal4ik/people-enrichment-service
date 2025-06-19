package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/adal4ik/people-enrichment-service/internal/models"
	"github.com/adal4ik/people-enrichment-service/internal/service"
	"github.com/go-chi/chi/v5"
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
	p.logger.Info("person inserted successfully", zap.String("name", person.Name))
	resp := models.APIResponse{
		Code:    201,
		Message: "Successfully created",
	}
	resp.Send(w)
}

func (p *PersonHandler) GetPersons(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	limitStr := query.Get("limit")
	offsetStr := query.Get("offset")
	name := query.Get("name")
	surname := query.Get("surname")
	gender := query.Get("gender")
	nationality := query.Get("nationality")

	ageMinStr := query.Get("age_min")
	ageMaxStr := query.Get("age_max")

	limit := 10
	offset := 0
	ageMin := 0
	ageMax := 200

	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}
	if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
		offset = o
	}

	if aMin, err := strconv.Atoi(ageMinStr); err == nil && aMin >= 0 {
		ageMin = aMin
	}
	if aMax, err := strconv.Atoi(ageMaxStr); err == nil && aMax >= 0 {
		ageMax = aMax
	}
	if ageMin > ageMax {
		p.handleError(w, req, 400, "age_min cannot be greater than age_max", nil)
		return
	}
	p.logger.Debug("GetPersons request params",
		zap.Int("limit", limit),
		zap.Int("offset", offset),
		zap.Int("age_min", ageMin),
		zap.Int("age_max", ageMax),
		zap.String("name", name),
		zap.String("surname", surname),
		zap.String("gender", gender),
		zap.String("nationality", nationality),
	)

	persons, err := p.service.GetPersons(req.Context(), limit, offset, ageMin, ageMax, name, surname, gender, nationality)
	if err != nil {
		p.handleError(w, req, 500, "failed to retrieve persons", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(persons); err != nil {
		p.handleError(w, req, 500, "failed to encode response", err)
		return
	}
	p.logger.Info("retrieved persons",
		zap.Int("count", len(persons)),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
		zap.String("name", name),
		zap.String("surname", surname))
}

func (p *PersonHandler) DeletePerson(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id == "" {
		p.handleError(w, req, 400, "id parameter is required", nil)
		return
	}

	err := p.service.DeletePerson(req.Context(), id)
	if err != nil {
		p.handleError(w, req, 500, "failed to delete person", err)
		return
	}
	p.logger.Info("person deleted successfully", zap.String("id", id))
	resp := models.APIResponse{
		Code:    200,
		Message: "Successfully deleted",
	}
	resp.Send(w)
}
