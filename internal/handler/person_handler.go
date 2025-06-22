package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/adal4ik/people-enrichment-service/internal/models"
	"github.com/adal4ik/people-enrichment-service/internal/service"
	"github.com/adal4ik/people-enrichment-service/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

	jsonErr := utils.APIError{
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
	if person.Name == "" {
		p.handleError(w, req, 400, "name is required", nil)
		return
	}
	if person.Surname == "" {
		p.handleError(w, req, 400, "surname is required", nil)
		return
	}
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
	resp := utils.APIResponse{
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

func (p *PersonHandler) GetPerson(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	uuidValue, err := uuid.Parse(id)
	if err != nil {
		p.handleError(w, req, 400, "invalid UUID format for id", err)
		return
	}
	p.logger.Debug("GetPerson request", zap.String("id", id))
	person, err := p.service.GetPerson(req.Context(), uuidValue)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			p.handleError(w, req, 404, "person not found", nil)
			return
		}
		p.handleError(w, req, 500, "failed to retrieve person", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(person); err != nil {
		p.handleError(w, req, 500, "failed to encode response", err)
		return
	}
	p.logger.Info("retrieved person",
		zap.String("id", id),
		zap.String("name", person.Name),
		zap.String("surname", person.Surname),
		zap.Int("age", person.Age),
		zap.String("gender", person.Gender),
		zap.String("nationality", person.Nationality),
	)
}

func (p *PersonHandler) DeletePerson(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id == "" {
		p.handleError(w, req, 400, "id parameter is required", nil)
		return
	}
	p.logger.Debug("DeletePerson request", zap.String("id", id))
	uuidValue, err := uuid.Parse(id)
	if err != nil {
		p.handleError(w, req, 400, "invalid UUID format for id", err)
		return
	}
	err = p.service.DeletePerson(req.Context(), uuidValue)
	if err != nil {
		p.handleError(w, req, 500, "failed to delete person", err)
		return
	}
	p.logger.Info("person deleted successfully", zap.String("id", id))
	resp := utils.APIResponse{
		Code:    200,
		Message: "Successfully deleted",
	}
	resp.Send(w)
}

func (p *PersonHandler) UpdatePerson(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	uuidValue, err := uuid.Parse(id)
	if err != nil {
		p.handleError(w, req, 400, "invalid UUID format for id", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var person models.Person
	err = decoder.Decode(&person)

	if err != nil {
		p.handleError(w, req, 400, "failed to decode request body", err)
		return
	}
	if person.Name == "" || person.Surname == "" {
		p.handleError(w, req, 400, "name and surname are required", nil)
		return
	}
	if person.Age == 0 || person.Gender == "" || person.Nationality == "" {
		p.handleError(w, req, 400, "age, gender, and nationality must not be empty", nil)
		return
	}
	person.ID = uuidValue

	err = p.service.UpdatePerson(req.Context(), person)
	if errors.Is(err, utils.ErrPersonNotFound) {
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}
	if err != nil {
		p.handleError(w, req, 500, "failed to update person", err)
		return
	}

	p.logger.Info("person updated successfully", zap.String("id", id))
	resp := utils.APIResponse{
		Code:    200,
		Message: "Successfully updated",
	}
	resp.Send(w)
}
