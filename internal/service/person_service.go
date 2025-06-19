package service

import "github.com/adal4ik/people-enrichment-service/internal/repository"

type PersonServiceInterface interface {
}

type PersonService struct {
	repo repository.PersonRepositoryInterface
}

func NewPersonService(repo repository.PersonRepositoryInterface) *PersonService {
	return &PersonService{
		repo: repo,
	}
}
