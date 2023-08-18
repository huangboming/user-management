package services

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"usermanagement/internal/models"
)

var DataFilePath = "../internal/services/data/users.json"

type UserServiceInterface interface {
	GetAllUsers() []models.User
	CreateUser(user models.User) error
	SearchUserByID(ID string) (models.User, error)
	SearchUserByUsername(username string) (models.User, error)
}

type UserService struct {
	Userdata []models.User
}

func NewUserService() (UserServiceInterface, error) {

	// import data from JSON file
	userdata := make([]models.User, 0)
	file, err := os.ReadFile(DataFilePath)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if len(file) == 0 {
		// if empty file
		return &UserService{
			Userdata: userdata,
		}, nil
	}

	err = json.Unmarshal([]byte(file), &userdata)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &UserService{
		Userdata: userdata,
	}, nil
}

// GetAllUsers get all users in the database and return a list of user
func (u *UserService) GetAllUsers() []models.User {
	return u.Userdata
}

// CreateUser adds a new user to the UserService.
// If the user with the same username already exists, it returns an error.
func (u *UserService) CreateUser(user models.User) error {

	// if the user already exists, return error
	_, err := u.SearchUserByUsername(user.Username)
	if err == nil {
		return errors.New("user already exsit")
	}

	// get all user data and append new user
	users := u.GetAllUsers()
	users = append(users, user)

	encoded, err := json.Marshal(users)
	if err != nil {
		log.Println(err)
		return err
	}

	// write data into file
	err = os.WriteFile(DataFilePath, encoded, 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	// everything's ok
	u.Userdata = users
	return nil
}

// SearchUserByID searches for a user in the database by the given ID.
// It returns the matched user and nil error if found, otherwise it returns an empty User model
// and an error indicating the user was not found.
func (u *UserService) SearchUserByID(ID string) (models.User, error) {
	users := u.GetAllUsers()
	for _, user := range users {
		if user.ID == ID {
			// find a user
			return user, nil
		}
	}
	return models.User{}, errors.New("not found")
}

// SearchUserByUsername searches for a user in the database by the given username.
// It returns the matched user and nil error if found, otherwise it returns an empty User model
// and an error indicating the user was not found.
func (u *UserService) SearchUserByUsername(username string) (models.User, error) {
	users := u.GetAllUsers()
	for _, user := range users {
		if user.Username == username {
			// find a user
			return user, nil
		}
	}
	return models.User{}, errors.New("not found")
}
