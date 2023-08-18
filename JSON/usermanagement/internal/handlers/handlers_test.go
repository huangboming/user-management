package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"usermanagement/internal/handlers"
	"usermanagement/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(user models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) GetAllUsers() []models.User {
	args := m.Called()
	return args.Get(0).([]models.User)
}

func (m *MockUserService) SearchUserByUsername(username string) (models.User, error) {
	args := m.Called(username)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserService) SearchUserByID(ID string) (models.User, error) {
	args := m.Called(ID)
	return args.Get(0).(models.User), args.Error(1)
}

func TestHandleRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		body       interface{}
		mockSetup  func(m *MockUserService)
		wantStatus int
	}{
		{
			// test case 1: successful registration, return http.StatusOK
			name: "successful registration",
			body: map[string]string{
				"username": "testuser",
				"password": "testpass",
			},
			mockSetup: func(m *MockUserService) {
				m.On("SearchUserByUsername", "testuser").Return(models.User{}, errors.New("not found"))
				m.On("CreateUser", mock.AnythingOfType("models.User")).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			// test case 2: user already exist, return http.StatusBadRequest
			name: "username already exists",
			body: map[string]string{
				"username": "existinguser",
				"password": "testpass",
			},
			mockSetup: func(m *MockUserService) {
				m.On("SearchUserByUsername", "existinguser").Return(models.User{Username: "existinguser"}, nil)
				m.On("CreateUser", mock.AnythingOfType("models.User")).Return(errors.New("user already exists"))
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			// test case 3: empty username, return http.StatusBadRequest
			name: "empty user name",
			body: map[string]string{
				"username": "",
				"password": "testpass",
			},
			mockSetup: func(m *MockUserService) {
				m.On("SearchUserByUsername", "").Return(models.User{}, nil)
				m.On("CreateUser", mock.AnythingOfType("models.User")).Return(errors.New("empty user name"))
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			// test case 4: invalid JSON, return http.StatusBadRequest
			name: "invalid JSON",
			body: "test body",
			mockSetup: func(m *MockUserService) {
				m.On("SearchUserByUsername", "").Return(models.User{}, nil)
				m.On("CreateUser", mock.AnythingOfType("models.User")).Return(errors.New("empty user name"))
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup mock
			mockUserService := new(MockUserService)
			tt.mockSetup(mockUserService)

			// setup router
			server := handlers.NewServer(mockUserService)
			server.SetupRoute()

			bodyBytes, _ := json.Marshal(tt.body)
			req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(bodyBytes))
			assert.NoError(t, err, "Should be able to create a request")

			resp := httptest.NewRecorder()
			server.GetRouter().ServeHTTP(resp, req)

			assert.Equal(t, tt.wantStatus, resp.Code, "Unexpected response status")
		})
	}
}

func TestHandleLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		body       interface{}
		mockSetup  func(m *MockUserService)
		wantStatus int
	}{
		{
			// test case 1: successful login, return http.StatusOK
			name: "successful login",
			body: map[string]string{
				"username": "testuser",
				"password": "testpass",
			},
			mockSetup: func(m *MockUserService) {
				m.On("SearchUserByUsername", "testuser").Return(models.User{
					Username: "testuser",
					Password: "testpass",
					ID:       "testid",
				}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			// test case 2: invalid password, return http.StatusBadRequest
			name: "invalid password",
			body: map[string]string{
				"username": "testuser",
				"password": "wrongpass",
			},
			mockSetup: func(m *MockUserService) {
				m.On("SearchUserByUsername", "testuser").Return(models.User{
					Username: "testuser",
					Password: "testpass",
					ID:       "testid",
				}, nil)
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			// test case 3: user not found, return http.StatusBadRequest
			name: "user not found",
			body: map[string]string{
				"username": "testuser",
				"password": "testpass",
			},
			mockSetup: func(m *MockUserService) {
				m.On("SearchUserByUsername", "testuser").Return(models.User{}, errors.New("not found"))
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			// test case 4: invalid JSON, return http.StatusBadRequest
			name: "invalid JSON",
			body: "test body",
			mockSetup: func(m *MockUserService) {
				m.On("SearchUserByUsername", "").Return(models.User{}, nil)
				m.On("CreateUser", mock.AnythingOfType("models.User")).Return(errors.New("empty user name"))
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MockUserService := new(MockUserService)
			tt.mockSetup(MockUserService)

			server := handlers.NewServer(MockUserService)
			server.SetupRoute()

			bodyBytes, _ := json.Marshal(tt.body)
			req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(bodyBytes))
			assert.NoError(t, err, "Should be able to create a request")

			resp := httptest.NewRecorder()
			server.GetRouter().ServeHTTP(resp, req)

			assert.Equal(t, tt.wantStatus, resp.Code, "Unexpected response status")
		})
	}
}

func TestHandleGetAllUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		mockSetup  func(m *MockUserService)
		wantStatus int
	}{
		{
			// test case 1: successful get all users, return http.StatusOK
			name: "successful get all users",
			mockSetup: func(m *MockUserService) {
				m.On("GetAllUsers").Return([]models.User{
					{
						Username: "testuser1",
						Password: "testpass1",
						ID:       "testid1",
					},
					{
						Username: "testuser2",
						Password: "testpass2",
						ID:       "testid2",
					},
				})
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MockUserService := new(MockUserService)
			tt.mockSetup(MockUserService)

			server := handlers.NewServer(MockUserService)
			server.SetupRoute()

			req, err := http.NewRequest(http.MethodGet, "/users", nil)
			assert.NoError(t, err, "Should be able to create a request")

			resp := httptest.NewRecorder()
			server.GetRouter().ServeHTTP(resp, req)

			assert.Equal(t, tt.wantStatus, resp.Code, "Unexpected response status")
		})
	}
}

func TestHandleSearchUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		query      string
		mockSetup  func(m *MockUserService)
		wantStatus int
	}{
		{
			// test case 1: successful search user by username, return http.StatusOK
			name:  "successful search user by username",
			query: "username=testuser",
			mockSetup: func(m *MockUserService) {
				m.On("SearchUserByUsername", "testuser").Return(models.User{
					Username: "testuser",
					Password: "testpass",
					ID:       "testid",
				}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			// test case 2: successful search user by id, return http.StatusOK
			name:  "successful search user by id",
			query: "id=testid",
			mockSetup: func(m *MockUserService) {
				m.On("SearchUserByID", "testid").Return(models.User{
					Username: "testuser",
					Password: "testpass",
					ID:       "testid",
				}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			// test case 3.1: user not found, return http.StatusBadRequest
			name:  "user not found",
			query: "username=testuser",
			mockSetup: func(m *MockUserService) {
				m.On("SearchUserByUsername", "testuser").Return(models.User{}, errors.New("not found"))
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			// test case 3.2: user not found, return http.StatusBadRequest
			name:  "user not found",
			query: "username=testid",
			mockSetup: func(m *MockUserService) {
				m.On("SearchUserByID", "testid").Return(models.User{}, errors.New("not found"))
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			// test case 4: invalid query, return http.StatusBadRequest
			name:       "invalid query",
			query:      "invalidquery",
			mockSetup:  func(m *MockUserService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MockUserService := new(MockUserService)
			tt.mockSetup(MockUserService)

			server := handlers.NewServer(MockUserService)
			server.SetupRoute()

			req, err := http.NewRequest(http.MethodGet, "/search?"+tt.query, nil)
			assert.NoError(t, err, "Should be able to create a request")

			resp := httptest.NewRecorder()
			server.GetRouter().ServeHTTP(resp, req)
		})
	}
}
