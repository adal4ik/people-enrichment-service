package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/adal4ik/people-enrichment-service/internal/models"
	"go.uber.org/zap"
)

type PersonRepositoryInterface interface {
	CreatePerson(ctx context.Context, person models.Person) error
}

type PersonRepository struct {
	logger *zap.Logger
	db     *sql.DB
}

func NewPersonRepository(db *sql.DB, logger *zap.Logger) *PersonRepository {
	return &PersonRepository{
		db:     db,
		logger: logger,
	}
}

func (p *PersonRepository) CreatePerson(ctx context.Context, person models.Person) error {
	query := `
		INSERT INTO persons (name, surname, patronymic, age, gender, nationality, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, now(), now())
	`
	p.logger.Debug("executing insert query", zap.String("query", query), zap.Any("person", person))

	_, err := p.db.ExecContext(ctx, query,
		person.Name,
		person.Surname,
		person.Patronymic,
		person.Age,
		person.Gender,
		person.Nationality,
	)

	if err != nil {
		return fmt.Errorf("failed to insert person: %w", err)
	}

	p.logger.Info("person inserted successfully", zap.String("name", person.Name))
	return nil
}
