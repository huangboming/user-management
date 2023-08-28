package test

import (
	"errors"
	"usermanagement/internal/models"
	"usermanagement/internal/services"

	"github.com/stretchr/testify/mock"
)

// ----- mock for handlers_test.go -----

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) LoginMongo() {
	m.Called()
}

func (m *MockUserService) CreateUser(user models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) GetAllUsers() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserService) SearchUserByUsername(username string) (models.User, error) {
	args := m.Called(username)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserService) SearchUserByID(ID string) (models.User, error) {
	args := m.Called(ID)
	return args.Get(0).(models.User), args.Error(1)
}

type mockMongoDB struct {
	mock.Mock
}

func (m *mockMongoDB) Create(item interface{}) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *mockMongoDB) Read(filter interface{}, result interface{}) ([]interface{}, error) {
	args := m.Called(filter, result)
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *mockMongoDB) Update(filter interface{}, update interface{}) error {
	args := m.Called(filter, update)
	return args.Error(0)
}

func (m *mockMongoDB) Delete(filter interface{}) error {
	args := m.Called(filter)
	return args.Error(0)
}

// ----- mock for user_test.go -----

type mockUserService struct {
	userservice *services.UserService
	mock.Mock
}

func (m *mockUserService) SearchUserByUsername(username string) (models.User, error) {
	args := m.Called(username)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *mockUserService) CreateUser(user models.User) error {
	// if the user already exists, return error
	_, err := m.SearchUserByUsername(user.Username)
	if err == nil {
		return errors.New("user already exist")
	}

	// insert to MongoDB
	err = m.userservice.Database.Create(user)
	if err != nil {
		return err
	}
	return nil
}
