package service

import (
	"github.com/adal4ik/people-enrichment-service/internal/repository"
	"go.uber.org/zap"
)

type Service struct {
	PersonService *PersonService
}

func New(repo *repository.Repository, logger *zap.Logger) *Service {
	return &Service{
		PersonService: NewPersonService(repo.PersonRepository, logger),
	}
}
