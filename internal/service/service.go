package service

import "github.com/adal4ik/people-enrichment-service/internal/repository"

type Service struct {
	PersonService *PersonService
}

func New(repo *repository.Repository) *Service {
	return &Service{
		PersonService: NewPersonService(repo.PersonRepository),
	}
}
