package repository

import "database/sql"

type Repository struct {
	PersonRepository *PersonRepository
}

func New(db *sql.DB) *Repository {
	return &Repository{
		PersonRepository: NewPersonRepository(db),
	}
}
