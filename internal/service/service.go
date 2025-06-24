package service

import (
	"github.com/adal4ik/people-enrichment-service/internal/config"
	"github.com/adal4ik/people-enrichment-service/internal/repository"
	"go.uber.org/zap"
)

type Service struct {
	PersonService *PersonService
}

func New(repo *repository.Repository, cfg config.Config, logger *zap.Logger) *Service {
	return &Service{
		PersonService: NewPersonService(repo.PersonRepository, cfg, logger),
	}
}
