package service

import (
	"context"
	"errors"
	"testing"

	"github.com/adal4ik/people-enrichment-service/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type mockPersonRepo struct {
	mock.Mock
}

func (m *mockPersonRepo) CreatePerson(ctx context.Context, person models.Person) error {
	args := m.Called(ctx, person)
	return args.Error(0)
}

func (m *mockPersonRepo) GetPersons(ctx context.Context, limit, offset, ageMin, ageMax int, name, surname, gender, nationality string) ([]models.Person, error) {
	args := m.Called(ctx, limit, offset, ageMin, ageMax, name, surname, gender, nationality)
	return args.Get(0).([]models.Person), args.Error(1)
}

func (m *mockPersonRepo) DeletePerson(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockPersonRepo) UpdatePerson(ctx context.Context, person models.Person) error {
	args := m.Called(ctx, person)
	return args.Error(0)
}

type mockPersonService struct {
	*PersonService
	age         int
	gender      string
	nationality string
	ageErr      error
	genderErr   error
	natErr      error
}

func (m *mockPersonService) GetAge(ctx context.Context, name string) (int, error) {
	return m.age, m.ageErr
}

func (m *mockPersonService) GetGender(ctx context.Context, name string) (string, error) {
	return m.gender, m.genderErr
}

func (m *mockPersonService) GetNationality(ctx context.Context, name string) (string, error) {
	return m.nationality, m.natErr
}

func (m *mockPersonService) CreatePerson(ctx context.Context, person models.Person) error {
	age, err := m.GetAge(ctx, person.Name)
	if err != nil {
		return errors.New("failed to enrich age: " + err.Error())
	}
	person.Age = age

	gender, err := m.GetGender(ctx, person.Name)
	if err != nil {
		return errors.New("failed to enrich gender: " + err.Error())
	}
	person.Gender = gender

	nationality, err := m.GetNationality(ctx, person.Name)
	if err != nil {
		return errors.New("failed to enrich nationality: " + err.Error())
	}
	person.Nationality = nationality

	if m.PersonService == nil || m.PersonService.repo == nil {
		return errors.New("repository is not initialized")
	}
	return m.PersonService.repo.CreatePerson(ctx, person)
}

func TestCreatePerson_Success(t *testing.T) {
	repo := new(mockPersonRepo)
	logger := zap.NewNop()
	person := models.Person{Name: "John"}

	base := &PersonService{repo: repo, logger: logger}
	svc := &mockPersonService{
		PersonService: base,
		age:           74,
		gender:        "male",
		nationality:   "NG",
	}

	repo.On("CreatePerson", mock.Anything, mock.MatchedBy(func(p models.Person) bool {
		assert.Equal(t, "John", p.Name)
		assert.Equal(t, 74, p.Age)
		assert.Equal(t, "male", p.Gender)
		assert.Equal(t, "NG", p.Nationality)
		return true
	})).Return(nil)

	err := svc.CreatePerson(context.Background(), person)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCreatePerson_AgeError(t *testing.T) {
	repo := new(mockPersonRepo)
	logger := zap.NewNop()
	person := models.Person{Name: "John"}

	base := &PersonService{repo: repo, logger: logger}

	svc := &mockPersonService{
		PersonService: base,
		ageErr:        errors.New("age error"),
	}

	err := svc.CreatePerson(context.Background(), person)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to enrich age")

	repo.AssertNotCalled(t, "CreatePerson", mock.Anything, mock.Anything)
}

func TestCreatePerson_GenderError(t *testing.T) {
	repo := new(mockPersonRepo)
	logger := zap.NewNop()
	person := models.Person{Name: "John"}

	repo.On("CreatePerson", mock.Anything, mock.Anything).Return(nil).Maybe()

	base := &PersonService{repo: repo, logger: logger}
	svc := &mockPersonService{
		PersonService: base,
		age:           30,
		genderErr:     errors.New("gender error"),
	}

	err := svc.CreatePerson(context.Background(), person)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to enrich gender")
	repo.AssertNotCalled(t, "CreatePerson", mock.Anything, mock.Anything)
}

func TestCreatePerson_NationalityError(t *testing.T) {
	repo := new(mockPersonRepo)
	logger := zap.NewNop()
	person := models.Person{Name: "John"}

	repo.On("CreatePerson", mock.Anything, mock.Anything).Return(nil).Maybe()

	base := &PersonService{repo: repo, logger: logger}
	svc := &mockPersonService{
		PersonService: base,
		age:           30,
		gender:        "male",
		natErr:        errors.New("nat error"),
	}

	err := svc.CreatePerson(context.Background(), person)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to enrich nationality")
	repo.AssertNotCalled(t, "CreatePerson", mock.Anything, mock.Anything)
}

func TestGetPersons(t *testing.T) {
	repo := new(mockPersonRepo)
	logger := zap.NewNop()
	expected := []models.Person{
		{Name: "Alice"},
		{Name: "Bob"},
	}
	repo.On("GetPersons", mock.Anything, 10, 0, 0, 100, "a", "b", "f", "US").Return(expected, nil)

	svc := &PersonService{repo: repo, logger: logger}
	res, err := svc.GetPersons(context.Background(), 10, 0, 0, 100, "a", "b", "f", "US")
	assert.NoError(t, err)
	assert.Equal(t, expected, res)
	repo.AssertExpectations(t)
}

func TestGetPersons_Pagination(t *testing.T) {
	repo := new(mockPersonRepo)
	logger := zap.NewNop()
	expected := []models.Person{
		{Name: "Alice"},
		{Name: "Bob"},
	}
	limit := 2
	offset := 5
	repo.On("GetPersons", mock.Anything, limit, offset, 0, 100, "", "", "", "").Return(expected, nil)

	svc := &PersonService{repo: repo, logger: logger}
	res, err := svc.GetPersons(context.Background(), limit, offset, 0, 100, "", "", "", "")
	assert.NoError(t, err)
	assert.Equal(t, expected, res)
	repo.AssertExpectations(t)
}

func TestGetPersons_FilterByName(t *testing.T) {
	repo := new(mockPersonRepo)
	logger := zap.NewNop()
	expected := []models.Person{
		{Name: "Charlie"},
	}
	repo.On("GetPersons", mock.Anything, 10, 0, 0, 100, "Charlie", "", "", "").Return(expected, nil)

	svc := &PersonService{repo: repo, logger: logger}
	res, err := svc.GetPersons(context.Background(), 10, 0, 0, 100, "Charlie", "", "", "")
	assert.NoError(t, err)
	assert.Equal(t, expected, res)
	repo.AssertExpectations(t)
}

func TestGetPersons_FilterBySurname(t *testing.T) {
	repo := new(mockPersonRepo)
	logger := zap.NewNop()
	expected := []models.Person{
		{Name: "David", Surname: "Smith"},
	}
	repo.On("GetPersons", mock.Anything, 10, 0, 0, 100, "", "Smith", "", "").Return(expected, nil)

	svc := &PersonService{repo: repo, logger: logger}
	res, err := svc.GetPersons(context.Background(), 10, 0, 0, 100, "", "Smith", "", "")
	assert.NoError(t, err)
	assert.Equal(t, expected, res)
	repo.AssertExpectations(t)
}

func TestGetPersons_FilterByGender(t *testing.T) {
	repo := new(mockPersonRepo)
	logger := zap.NewNop()
	expected := []models.Person{
		{Name: "Eve", Gender: "female"},
	}
	repo.On("GetPersons", mock.Anything, 10, 0, 0, 100, "", "", "female", "").Return(expected, nil)

	svc := &PersonService{repo: repo, logger: logger}
	res, err := svc.GetPersons(context.Background(), 10, 0, 0, 100, "", "", "female", "")
	assert.NoError(t, err)
	assert.Equal(t, expected, res)
	repo.AssertExpectations(t)
}

func TestGetPersons_FilterByNationality(t *testing.T) {
	repo := new(mockPersonRepo)
	logger := zap.NewNop()
	expected := []models.Person{
		{Name: "Frank", Nationality: "DE"},
	}
	repo.On("GetPersons", mock.Anything, 10, 0, 0, 100, "", "", "", "DE").Return(expected, nil)

	svc := &PersonService{repo: repo, logger: logger}
	res, err := svc.GetPersons(context.Background(), 10, 0, 0, 100, "", "", "", "DE")
	assert.NoError(t, err)
	assert.Equal(t, expected, res)
	repo.AssertExpectations(t)
}

func TestGetPersons_FilterByAgeRange(t *testing.T) {
	repo := new(mockPersonRepo)
	logger := zap.NewNop()
	expected := []models.Person{
		{Name: "Grace", Age: 30},
	}
	repo.On("GetPersons", mock.Anything, 10, 0, 25, 35, "", "", "", "").Return(expected, nil)

	svc := &PersonService{repo: repo, logger: logger}
	res, err := svc.GetPersons(context.Background(), 10, 0, 25, 35, "", "", "", "")
	assert.NoError(t, err)
	assert.Equal(t, expected, res)
	repo.AssertExpectations(t)
}

func TestDeletePerson_Success(t *testing.T) {
	repo := new(mockPersonRepo)
	logger := zap.NewNop()
	svc := &PersonService{repo: repo, logger: logger}
	id := uuid.New()

	repo.On("DeletePerson", mock.Anything, id).Return(nil)

	err := svc.DeletePerson(context.Background(), id)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestDeletePerson_Error(t *testing.T) {
	repo := new(mockPersonRepo)
	logger := zap.NewNop()
	svc := &PersonService{repo: repo, logger: logger}
	id := uuid.New()
	expectedErr := errors.New("delete error")

	repo.On("DeletePerson", mock.Anything, id).Return(expectedErr)

	err := svc.DeletePerson(context.Background(), id)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	repo.AssertExpectations(t)
}

func TestUpdatePerson_Success(t *testing.T) {
	repo := new(mockPersonRepo)
	logger := zap.NewNop()
	svc := &PersonService{repo: repo, logger: logger}
	person := models.Person{Name: "Jane", Age: 25}

	repo.On("UpdatePerson", mock.Anything, person).Return(nil)

	err := svc.UpdatePerson(context.Background(), person)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestUpdatePerson_Error(t *testing.T) {
	repo := new(mockPersonRepo)
	logger := zap.NewNop()
	svc := &PersonService{repo: repo, logger: logger}
	person := models.Person{Name: "Jane", Age: 25}
	expectedErr := errors.New("update error")

	repo.On("UpdatePerson", mock.Anything, person).Return(expectedErr)

	err := svc.UpdatePerson(context.Background(), person)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	repo.AssertExpectations(t)
}

func (m *mockPersonRepo) GetPerson(ctx context.Context, id uuid.UUID) (models.Person, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Person), args.Error(1)
}

func TestGetPerson_Success(t *testing.T) {
	repo := new(mockPersonRepo)
	logger := zap.NewNop()
	svc := &PersonService{repo: repo, logger: logger}
	id := uuid.New()
	expected := models.Person{Name: "Test", Age: 42}

	repo.On("GetPerson", mock.Anything, id).Return(expected, nil)

	person, err := svc.GetPerson(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, expected, person)
	repo.AssertExpectations(t)
}

func TestGetPerson_Error(t *testing.T) {
	repo := new(mockPersonRepo)
	logger := zap.NewNop()
	svc := &PersonService{repo: repo, logger: logger}
	id := uuid.New()
	expectedErr := errors.New("not found")

	repo.On("GetPerson", mock.Anything, id).Return(models.Person{}, expectedErr)

	person, err := svc.GetPerson(context.Background(), id)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, models.Person{}, person)
	repo.AssertExpectations(t)
}
