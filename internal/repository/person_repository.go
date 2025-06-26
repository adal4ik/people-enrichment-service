package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/adal4ik/people-enrichment-service/internal/models"
	"github.com/adal4ik/people-enrichment-service/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PersonRepositoryInterface interface {
	CreatePerson(ctx context.Context, person models.Person) error
	GetPersons(ctx context.Context, limit, offset, ageMin, ageMax int, name, surname, gender, nationality string) ([]models.Person, error)
	GetPerson(ctx context.Context, id uuid.UUID) (models.Person, error)
	DeletePerson(ctx context.Context, id uuid.UUID) error
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

func (p *PersonRepository) GetPerson(ctx context.Context, id uuid.UUID) (models.Person, error) {
	query := `
		SELECT id, name, surname, patronymic, age, gender, nationality, created_at, updated_at
		FROM persons
		WHERE id = $1
	`
	p.logger.Debug("executing select query", zap.String("query", query), zap.Any("id", id))

	var person models.Person
	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&person.ID,
		&person.Name,
		&person.Surname,
		&person.Patronymic,
		&person.Age,
		&person.Gender,
		&person.Nationality,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Person{}, sql.ErrNoRows
		}
		return models.Person{}, fmt.Errorf("failed to get person: %w", err)
	}
	return person, nil
}

func (p *PersonRepository) DeletePerson(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM persons WHERE id = $1
	`
	p.logger.Debug("executing delete query", zap.String("query", query), zap.Any("id", id))

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
	query := "UPDATE persons SET "
	args := []interface{}{}
	i := 1

	query += fmt.Sprintf("name = $%d, ", i)
	args = append(args, person.Name)
	i++
	query += fmt.Sprintf("surname = $%d, ", i)
	args = append(args, person.Surname)
	i++

	if person.Patronymic != nil {
		query += fmt.Sprintf("patronymic = $%d, ", i)
		args = append(args, *person.Patronymic)
		i++
	}
	if person.Age != nil {
		query += fmt.Sprintf("age = $%d, ", i)
		args = append(args, *person.Age)
		i++
	}
	if person.Gender != nil {
		query += fmt.Sprintf("gender = $%d, ", i)
		args = append(args, *person.Gender)
		i++
	}
	if person.Nationality != nil {
		query += fmt.Sprintf("nationality = $%d, ", i)
		args = append(args, *person.Nationality)
		i++
	}

	query = strings.TrimSuffix(query, ", ")
	query += fmt.Sprintf(", updated_at = now() WHERE id = $%d", i)
	args = append(args, person.ID)

	p.logger.Debug("executing dynamic update query", zap.String("query", query), zap.Any("args", args))

	res, err := p.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update person: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return utils.ErrPersonNotFound
	}

	return nil
}
