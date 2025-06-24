package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adal4ik/people-enrichment-service/internal/config"
	"github.com/adal4ik/people-enrichment-service/internal/models"
	"github.com/adal4ik/people-enrichment-service/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PersonServiceInterface interface {
	CreatePerson(ctx context.Context, person models.Person) error
	GetPersons(ctx context.Context, limit, offset, ageMin, ageMax int, name, surname, gender, nationality string) ([]models.Person, error)
	GetPerson(ctx context.Context, id uuid.UUID) (models.Person, error)
	DeletePerson(ctx context.Context, id uuid.UUID) error
	UpdatePerson(ctx context.Context, person models.Person) error
	GetAge(ctx context.Context, name string) (int, error)
	GetGender(ctx context.Context, name string) (string, error)
	GetNationality(ctx context.Context, name string) (string, error)
}

type PersonService struct {
	repo   repository.PersonRepositoryInterface
	cfg    config.Config
	logger *zap.Logger
}

func NewPersonService(repo repository.PersonRepositoryInterface, cfg config.Config, logger *zap.Logger) *PersonService {
	return &PersonService{
		repo:   repo,
		cfg:    cfg,
		logger: logger,
	}
}

func (p *PersonService) CreatePerson(ctx context.Context, person models.Person) error {
	age, err := p.GetAge(ctx, person.Name)
	if err != nil {
		p.logger.Error("failed to get age", zap.Error(err))
		return fmt.Errorf("failed to enrich age: %w", err)
	}
	gender, err := p.GetGender(ctx, person.Name)
	if err != nil {
		p.logger.Error("failed to get gender", zap.Error(err))
		return fmt.Errorf("failed to enrich gender: %w", err)
	}
	nationality, err := p.GetNationality(ctx, person.Name)
	if err != nil {
		p.logger.Error("failed to get nationality", zap.Error(err))
		return fmt.Errorf("failed to enrich nationality: %w", err)
	}
	p.logger.Debug("enriched person data",
		zap.Int("age", age),
		zap.String("gender", gender),
		zap.String("nationality", nationality),
	)

	person.Age = age
	person.Gender = gender
	person.Nationality = nationality

	return p.repo.CreatePerson(ctx, person)
}

func (p *PersonService) GetPersons(ctx context.Context, limit, offset, age_min, age_max int, name, surname, gender, nationality string) ([]models.Person, error) {
	p.logger.Debug("Service GetPersons called",
		zap.Int("limit", limit),
		zap.Int("offset", offset),
		zap.Int("age_min", age_min),
		zap.Int("age_max", age_max),
		zap.String("name", name),
		zap.String("surname", surname),
		zap.String("gender", gender),
		zap.String("nationality", nationality),
	)
	return p.repo.GetPersons(ctx, limit, offset, age_min, age_max, name, surname, gender, nationality)
}

func (p *PersonService) GetPerson(ctx context.Context, id uuid.UUID) (models.Person, error) {
	p.logger.Debug("getting person by id", zap.Any("id", id))
	return p.repo.GetPerson(ctx, id)
}

func (p *PersonService) DeletePerson(ctx context.Context, id uuid.UUID) error {
	p.logger.Debug("deleting person", zap.Any("id", id))
	return p.repo.DeletePerson(ctx, id)
}

func (p *PersonService) UpdatePerson(ctx context.Context, person models.Person) error {
	p.logger.Debug("updating person", zap.Any("person", person))
	return p.repo.UpdatePerson(ctx, person)
}

func (p *PersonService) GetAge(ctx context.Context, name string) (int, error) {
	url := fmt.Sprintf("%s&name=%s", p.cfg.APIAgeURL, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to call agify API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("agify API returned status %d", resp.StatusCode)
	}
	var result struct {
		Age int `json:"age"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode agify API response: %w", err)
	}

	return result.Age, nil
}

func (p *PersonService) GetGender(ctx context.Context, name string) (string, error) {
	url := fmt.Sprintf("%s&name=%s", p.cfg.APIGenderURL, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call agify API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("agify API returned status %d", resp.StatusCode)
	}

	var result struct {
		Gender string `json:"gender"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode agify API response: %w", err)
	}
	return result.Gender, nil
}

func (p *PersonService) GetNationality(ctx context.Context, name string) (string, error) {
	url := fmt.Sprintf("%s&name=%s", p.cfg.APINationURL, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call nationalize API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("nationalize API returned status %d", resp.StatusCode)
	}

	var result struct {
		Country []struct {
			CountryID   string  `json:"country_id"`
			Probability float64 `json:"probability"`
		} `json:"country"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode nationalize API response: %w", err)
	}

	if len(result.Country) > 0 {
		return result.Country[0].CountryID, nil
	}

	return "", nil
}
