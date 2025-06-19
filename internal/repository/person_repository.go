package repository

import "database/sql"

type PersonRepositoryInterface interface {
}

type PersonRepository struct {
	db *sql.DB
}

func NewPersonRepository(db *sql.DB) *PersonRepository {
	return &PersonRepository{db: db}
}
