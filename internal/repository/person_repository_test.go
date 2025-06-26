package repository

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/adal4ik/people-enrichment-service/internal/models"
	"github.com/adal4ik/people-enrichment-service/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func newTestRepo(t *testing.T) (*PersonRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	logger := zaptest.NewLogger(t)
	repo := NewPersonRepository(db, logger)
	return repo, mock, func() { db.Close() }
}

func TestCreatePerson_Success(t *testing.T) {
	repo, mock, close := newTestRepo(t)
	defer close()

	person := models.Person{
		Name:        "John",
		Surname:     "Doe",
		Patronymic:  nil,
		Age:         nil,
		Gender:      nil,
		Nationality: nil,
	}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO persons (name, surname, patronymic, age, gender, nationality, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, now(), now())")).
		WithArgs(person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.CreatePerson(context.Background(), person)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePerson_Error(t *testing.T) {
	repo, mock, close := newTestRepo(t)
	defer close()

	person := models.Person{Name: "Jane", Surname: "Smith"}
	mock.ExpectExec("INSERT INTO persons").
		WillReturnError(errors.New("insert error"))

	err := repo.CreatePerson(context.Background(), person)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to insert person")
}

func TestGetPersons_QueryError(t *testing.T) {
	repo, mock, close := newTestRepo(t)
	defer close()

	mock.ExpectQuery("SELECT id, name, surname, patronymic, age, gender, nationality, created_at, updated_at FROM persons").
		WillReturnError(errors.New("query error"))

	_, err := repo.GetPersons(context.Background(), 10, 0, 0, 100, "", "", "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to query persons")
}

func TestGetPersons_ScanError(t *testing.T) {
	repo, mock, close := newTestRepo(t)
	defer close()

	rows := sqlmock.NewRows([]string{"id", "name", "surname", "patronymic", "age", "gender", "nationality", "created_at", "updated_at"}).
		AddRow("bad-uuid", "John", "Doe", nil, 30, "male", "US", time.Now(), time.Now())

	mock.ExpectQuery("SELECT id, name, surname, patronymic, age, gender, nationality, created_at, updated_at FROM persons").
		WillReturnRows(rows)

	_, err := repo.GetPersons(context.Background(), 10, 0, 0, 100, "", "", "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to scan person")
}

func TestGetPerson_NotFound(t *testing.T) {
	repo, mock, close := newTestRepo(t)
	defer close()

	id := uuid.New()
	mock.ExpectQuery("SELECT id, name, surname, patronymic, age, gender, nationality, created_at, updated_at FROM persons WHERE id = \\$1").
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetPerson(context.Background(), id)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestGetPerson_ScanError(t *testing.T) {
	repo, mock, close := newTestRepo(t)
	defer close()

	id := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "name", "surname", "patronymic", "age", "gender", "nationality", "created_at", "updated_at"}).
		AddRow("bad-uuid", "John", "Doe", nil, 30, "male", "US", time.Now(), time.Now())

	mock.ExpectQuery("SELECT id, name, surname, patronymic, age, gender, nationality, created_at, updated_at FROM persons WHERE id = \\$1").
		WithArgs(id).
		WillReturnRows(rows)

	_, err := repo.GetPerson(context.Background(), id)
	assert.Error(t, err)
}

func TestDeletePerson_Success(t *testing.T) {
	repo, mock, close := newTestRepo(t)
	defer close()

	id := uuid.New()
	mock.ExpectExec("DELETE FROM persons WHERE id = \\$1").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeletePerson(context.Background(), id)
	assert.NoError(t, err)
}

func TestDeletePerson_NotFound(t *testing.T) {
	repo, mock, close := newTestRepo(t)
	defer close()

	id := uuid.New()
	mock.ExpectExec("DELETE FROM persons WHERE id = \\$1").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.DeletePerson(context.Background(), id)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestDeletePerson_ExecError(t *testing.T) {
	repo, mock, close := newTestRepo(t)
	defer close()

	id := uuid.New()
	mock.ExpectExec("DELETE FROM persons WHERE id = \\$1").
		WithArgs(id).
		WillReturnError(errors.New("delete error"))

	err := repo.DeletePerson(context.Background(), id)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete person")
}

func TestUpdatePerson_Success(t *testing.T) {
	repo, mock, close := newTestRepo(t)
	defer close()

	id := uuid.New()
	age := 25
	gender := "male"
	person := models.Person{
		ID:         id,
		Name:       "John",
		Surname:    "Doe",
		Patronymic: nil,
		Age:        &age,
		Gender:     &gender,
	}

	mock.ExpectExec("UPDATE persons SET name = \\$1, surname = \\$2, age = \\$3, gender = \\$4, updated_at = now\\(\\) WHERE id = \\$5").
		WithArgs(person.Name, person.Surname, *person.Age, *person.Gender, person.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdatePerson(context.Background(), person)
	assert.NoError(t, err)
}

func TestUpdatePerson_NotFound(t *testing.T) {
	repo, mock, close := newTestRepo(t)
	defer close()

	id := uuid.New()
	age := 25
	person := models.Person{
		ID:      id,
		Name:    "John",
		Surname: "Doe",
		Age:     &age,
	}

	mock.ExpectExec("UPDATE persons SET name = \\$1, surname = \\$2, age = \\$3, updated_at = now\\(\\) WHERE id = \\$4").
		WithArgs(person.Name, person.Surname, *person.Age, person.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.UpdatePerson(context.Background(), person)
	assert.ErrorIs(t, err, utils.ErrPersonNotFound)
}

func TestUpdatePerson_ExecError(t *testing.T) {
	repo, mock, close := newTestRepo(t)
	defer close()

	id := uuid.New()
	age := 25
	person := models.Person{
		ID:      id,
		Name:    "John",
		Surname: "Doe",
		Age:     &age,
	}

	mock.ExpectExec("UPDATE persons SET name = \\$1, surname = \\$2, age = \\$3, updated_at = now\\(\\) WHERE id = \\$4").
		WithArgs(person.Name, person.Surname, *person.Age, person.ID).
		WillReturnError(errors.New("update error"))

	err := repo.UpdatePerson(context.Background(), person)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update person")
}
