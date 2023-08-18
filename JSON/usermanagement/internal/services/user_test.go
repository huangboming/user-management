package services_test

import (
	"log"
	"os"
	"testing"
	"usermanagement/internal/models"
	"usermanagement/internal/services"

	"github.com/stretchr/testify/assert"
)

// setupMockData creates a new UserService with mock data
func setupMockData() *services.UserService {
	// Mock data for testing
	users := []models.User{
		{
			ID:       "1",
			Username: "testuser1",
			Password: "testpass1",
		},
		{
			ID:       "2",
			Username: "testuser2",
			Password: "testpass2",
		},
		{
			ID:       "3",
			Username: "testuser3",
			Password: "testpass3",
		},
	}

	return &services.UserService{
		Userdata: users,
	}
}

func TestGetAllUsers(t *testing.T) {
	service := setupMockData()
	users := service.GetAllUsers()

	assert.Equal(t, 3, len(users), "Expected 3 users")
}

func TestCreateUser(t *testing.T) {

	service := setupMockData()

	// 因为在CreateUser中有一段代码是从文件中读取数据，所以这里需要mock一个临时文件
	temfile, err := os.CreateTemp("", "users.json")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(temfile.Name())

	services.DataFilePath = temfile.Name()

	// Test creating a new user
	newUser := models.User{
		ID:       "4",
		Username: "testuser4",
		Password: "testpass4",
	}
	err = service.CreateUser(newUser)
	assert.Nil(t, err, "Expected no error when creating a new user")

	// Test creating a user with an existing username
	err = service.CreateUser(newUser)
	assert.NotNil(t, err, "Expected an error when creating a user with an existing username")
}

func TestSearchUserByID(t *testing.T) {
	service := setupMockData()

	// Test searching for an existing user by ID
	user, err := service.SearchUserByID("1")
	assert.Nil(t, err, "Expected no error when searching for an existing user by ID")
	assert.Equal(t, "testuser1", user.Username, "Expected to find user with username 'testuser1'")

	// Test searching for a non-existing user by ID
	_, err = service.SearchUserByID("10")
	assert.NotNil(t, err, "Expected an error when searching for a non-existing user by ID")
}

func TestSearchUserByUsername(t *testing.T) {
	service := setupMockData()

	// Test searching for an existing user by username
	user, err := service.SearchUserByUsername("testuser1")
	assert.Nil(t, err, "Expected no error when searching for an existing user by username")
	assert.Equal(t, "1", user.ID, "Expected to find user with ID '1'")

	// Test searching for a non-existing user by username
	_, err = service.SearchUserByUsername("unknownuser")
	assert.NotNil(t, err, "Expected an error when searching for a non-existing user by username")
}
