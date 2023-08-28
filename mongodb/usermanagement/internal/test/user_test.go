package test

import (
	"errors"
	"testing"
	"usermanagement/internal/models"
	"usermanagement/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestGetAllUsers tests the GetAllUsers method of the UserService
func TestGetAllUsers(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(m *MockDB)
		wantError bool
		wantUsers int
	}{
		{
			// test case 1: successfully get all users
			name: "successfully get all users",
			mockSetup: func(m *MockDB) {
				m.On("Read", mock.Anything, mock.Anything).Return([]interface{}{
					&models.User{
						ID:       primitive.NewObjectID(),
						Username: "testuser1",
						Password: "testpass1",
					},
					&models.User{
						ID:       primitive.NewObjectID(),
						Username: "testuser2",
						Password: "testpass2",
					},
				}, nil)
			},
			wantError: false,
			wantUsers: 2,
		},
		{
			// test case 2: failed to get all users
			name: "failed to get all users",
			mockSetup: func(m *MockDB) {
				m.On("Read", mock.Anything, mock.Anything).Return([]interface{}{}, errors.New("failed to read from db"))
			},
			wantError: true,
			wantUsers: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.mockSetup(mockDB)

			// create a new UserService with the mockDB
			userService := services.NewUserService()
			userService.Database = mockDB

			// call the GetAllUsers method
			users, err := userService.GetAllUsers()

			if tt.wantError && err == nil {
				t.Fatal("expected an error but got none")
			}

			if !tt.wantError && err != nil {
				t.Fatalf("did not expect an error but got: %v", err)
			}

			// assert the returned users
			if len(users) != tt.wantUsers {
				t.Errorf("expected %v users, got %v", tt.wantUsers, len(users))
			}

			// assert the mockDB
			mockDB.AssertExpectations(t)
		})
	}
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name      string
		inputUser models.User
		mockSetup func(m *mockUserService, db *MockDB)
		wantErr   bool
	}{
		{
			name:      "successfully create user",
			inputUser: models.User{Username: "testuser", Password: "testpass"},
			mockSetup: func(m *mockUserService, db *MockDB) {
				m.On("SearchUserByUsername", "testuser").Return(models.User{}, errors.New("user not found"))
				db.On("Create", mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "user already exists",
			inputUser: models.User{Username: "testuser", Password: "testpass"},
			mockSetup: func(m *mockUserService, db *MockDB) {
				m.On("SearchUserByUsername", "testuser").Return(models.User{
					Username: "testuser",
					Password: "testpass",
					ID:       primitive.NewObjectID(),
				}, nil)
			},
			wantErr: true,
		},
		{
			name:      "database error on create",
			inputUser: models.User{Username: "testuser", Password: "testpass"},
			mockSetup: func(m *mockUserService, db *MockDB) {
				m.On("SearchUserByUsername", "testuser").Return(models.User{}, errors.New("user not found"))
				db.On("Create", mock.Anything).Return(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			mockService := &mockUserService{userservice: &services.UserService{Database: mockDB}}
			tt.mockSetup(mockService, mockDB)

			err := mockService.CreateUser(tt.inputUser)

			if tt.wantErr {
				assert.Error(t, err, "Expected an error")
			} else {
				assert.NoError(t, err, "Did not expect an error")
			}

			mockService.AssertExpectations(t)
			mockDB.AssertExpectations(t)
		})
	}
}

// TestSearchUserByID tests the SearchUserByID method of the UserService
func TestSearchUserByID(t *testing.T) {

	tests := []struct {
		name      string
		ID        string
		mockSetup func(m *MockDB)
		wantError bool
		wantUser  bool
	}{
		{
			// test case 1: successfully get user by ID
			name: "successfully get user by ID",
			ID:   primitive.NewObjectID().Hex(),
			mockSetup: func(m *MockDB) {
				m.On("Read", mock.Anything, mock.Anything).Return([]interface{}{
					&models.User{
						ID:       primitive.NewObjectID(),
						Username: "testuser1",
						Password: "testpass1",
					},
				}, nil)
			},
			wantError: false,
			wantUser:  true,
		},
		{
			// test case 2: failed to get user by ID
			name: "failed to get user by ID",
			ID:   primitive.NewObjectID().Hex(),
			mockSetup: func(m *MockDB) {
				m.On("Read", mock.Anything, mock.Anything).Return([]interface{}{}, errors.New("failed to read from db"))
			},
			wantError: true,
			wantUser:  false,
		},
		{
			// test case 3: user not found
			name: "user not found",
			ID:   primitive.NewObjectID().Hex(),
			mockSetup: func(m *MockDB) {
				m.On("Read", mock.Anything, mock.Anything).Return([]interface{}{}, nil)
			},
			wantError: true,
			wantUser:  false,
		},
		{
			// test case 4: error when converting ID to ObjectID
			name: "error when converting ID to ObjectID",
			ID:   "invalidID",
			mockSetup: func(m *MockDB) {
				m.On("Read", mock.Anything, mock.Anything).Return([]interface{}{}, nil)
			},
			wantError: true,
			wantUser:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.mockSetup(mockDB)

			// create a new UserService with the mockDB
			userService := services.NewUserService()
			userService.Database = mockDB

			// call the SearchUserByID method
			user, err := userService.SearchUserByID(tt.ID)

			if tt.wantError {
				assert.Error(t, err, "Expected an error")
			} else {
				assert.NoError(t, err, "Did not expect an error")
			}

			if tt.wantUser {
				assert.NotEmpty(t, user.Username, "Expected a user but got none")
			} else {
				assert.Empty(t, user.Username, "Did not expect a user")
			}
		})
	}
}

// TestSearchUserByUsername tests the SearchUserByUsername method of the UserService
func TestSearchUserByUsername(t *testing.T) {

	tests := []struct {
		name      string
		mockSetup func(m *MockDB)
		wantError bool
		wantUser  bool
	}{
		{
			// test case 1: successfully get user by username
			name: "successfully get user by username",
			mockSetup: func(m *MockDB) {
				m.On("Read", mock.Anything, mock.Anything).Return([]interface{}{
					&models.User{
						ID:       primitive.NewObjectID(),
						Username: "testuser1",
						Password: "testpass1",
					},
				}, nil)
			},
			wantError: false,
			wantUser:  true,
		},
		{
			// test case 2: failed to get user by username
			name: "failed to get user by username",
			mockSetup: func(m *MockDB) {
				m.On("Read", mock.Anything, mock.Anything).Return([]interface{}{}, errors.New("failed to read from db"))
			},
			wantError: true,
			wantUser:  false,
		},
		{
			// test case 3: user not found
			name: "user not found",
			mockSetup: func(m *MockDB) {
				m.On("Read", mock.Anything, mock.Anything).Return([]interface{}{}, nil)
			},
			wantError: true,
			wantUser:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.mockSetup(mockDB)

			// create a new UserService with the mockDB
			userService := services.NewUserService()
			userService.Database = mockDB

			// call the SearchUserByID method
			user, err := userService.SearchUserByUsername("testuser1")

			if tt.wantError && err == nil {
				t.Fatal("expected an error but got none")
			}

			if !tt.wantError && err != nil {
				t.Fatalf("did not expect an error but got: %v", err)
			}

			// assert the returned user
			if tt.wantUser && user.Username == "" {
				t.Errorf("expected a user but got none")
			}

			// assert the mockDB
			mockDB.AssertExpectations(t)
		})
	}
}
