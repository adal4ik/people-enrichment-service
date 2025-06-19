package repository

import (
	"database/sql"

	"go.uber.org/zap"
)

type Repository struct {
	PersonRepository *PersonRepository
}

func New(db *sql.DB, logger *zap.Logger) *Repository {
	return &Repository{
		PersonRepository: NewPersonRepository(db, logger),
	}
}
