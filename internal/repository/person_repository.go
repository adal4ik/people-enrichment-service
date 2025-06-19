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
	GetPersons(ctx context.Context, limit, offset, ageMin, ageMax int, name, surname, gender, nationality string) ([]models.Person, error)
	DeletePerson(ctx context.Context, id string) error
	UpdatePerson(ctx context.Context, person models.Person) error
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

	return nil
}

func (p *PersonRepository) GetPersons(ctx context.Context, limit, offset, ageMin, ageMax int, name, surname, gender, nationality string) ([]models.Person, error) {
	query := `
		SELECT id, name, surname, patronymic, age, gender, nationality, created_at, updated_at
		FROM persons
		WHERE age BETWEEN $1 AND $2
	`

	args := []interface{}{ageMin, ageMax}
	argPos := 3

	if name != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", argPos)
		args = append(args, "%"+name+"%")
		argPos++
	}
	if surname != "" {
		query += fmt.Sprintf(" AND surname ILIKE $%d", argPos)
		args = append(args, "%"+surname+"%")
		argPos++
	}
	if gender != "" {
		query += fmt.Sprintf(" AND gender = $%d", argPos)
		args = append(args, gender)
		argPos++
	}
	if nationality != "" {
		query += fmt.Sprintf(" AND nationality = $%d", argPos)
		args = append(args, nationality)
		argPos++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, limit, offset)

	p.logger.Debug("executing query", zap.String("query", query), zap.Any("args", args))

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query persons: %w", err)
	}
	defer rows.Close()

	var persons []models.Person
	for rows.Next() {
		var person models.Person
		if err := rows.Scan(
			&person.ID,
			&person.Name,
			&person.Surname,
			&person.Patronymic,
			&person.Age,
			&person.Gender,
			&person.Nationality,
			&person.CreatedAt,
			&person.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan person: %w", err)
		}
		persons = append(persons, person)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return persons, nil
}

func (p *PersonRepository) DeletePerson(ctx context.Context, id string) error {
	query := `
		DELETE FROM persons WHERE id = $1
	`
	p.logger.Debug("executing delete query", zap.String("query", query), zap.String("id", id))

	result, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete person: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (p *PersonRepository) UpdatePerson(ctx context.Context, person models.Person) error {
	query := `UPDATE persons SET name = $1, surname = $2, patronymic = $3, age = $4, gender = $5, nationality = $6, updated_at = now() WHERE id = $7`
	p.logger.Debug("executing update query", zap.String("query", query), zap.Any("person", person))
	_, err := p.db.ExecContext(ctx, query,
		person.Name,
		person.Surname,
		person.Patronymic,
		person.Age,
		person.Gender,
		person.Nationality,
		person.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update person: %w", err)
	}
	return nil
}
