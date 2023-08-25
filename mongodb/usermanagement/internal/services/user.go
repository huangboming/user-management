package services

import (
	"context"
	"errors"
	"log"
	"os"
	"usermanagement/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type UserService struct {
	Database models.CURDInterface
}

type UserServiceInterface interface {
	LoginMongo()
	GetAllUsers() ([]models.User, error)
	CreateUser(user models.User) error
	SearchUserByID(ID string) (models.User, error)
	SearchUserByUsername(username string) (models.User, error)
}

func NewUserService() *UserService {
	return &UserService{
		Database: nil,
	}
}

// LoginMongo: login mongodb
func (u *UserService) LoginMongo() {
	db := models.NewMongoDB()
	client, _ := mongo.Connect(db.Ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

	db.Client = client
	db.Collection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("user")

	u.Database = db
}

// GetAllUsers get all users in the database and return a list of user
func (u *UserService) GetAllUsers() ([]models.User, error) {
	// select data from mongodb
	found, err := u.Database.Read(bson.M{}, &models.User{})
	if err != nil {
		return nil, err
	}
	users := make([]models.User, 0)
	for _, user := range found {
		user, ok := user.(*models.User)
		if !ok {
			return nil, errors.New("type assertion failed")
		}
		users = append(users, *user)
	}
	return users, err
}

// CreateUser adds a new user to the UserService.
// If the user with the same username already exists, it returns an error.
func (u *UserService) CreateUser(user models.User) error {

	// if the user already exists, return error
	_, err := u.SearchUserByUsername(user.Username)
	if err == nil {
		return errors.New("user already exist")
	}

	// insert to MongoDB
	err = u.Database.Create(user)
	if err != nil {
		return err
	}
	return nil
}

// SearchUserByID searches for a user in the database by the given ID.
// It returns the matched user and nil error if found, otherwise it returns an empty User model
// and an error indicating the user was not found.
func (u *UserService) SearchUserByID(ID string) (models.User, error) {
	// convert id from string to MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return models.User{}, err
	}

	// search from MongoDB
	filter := bson.M{"_id": objectID}
	found, err := u.Database.Read(filter, &models.User{})
	if err != nil {
		return models.User{}, err
	}

	// found one
	if len(found) > 0 {
		user, ok := found[0].(*models.User)
		if !ok {
			return models.User{}, errors.New("type assertion failed")
		}
		return *user, nil
	}

	// not found
	return models.User{}, errors.New("not found")
}

// SearchUserByUsername searches for a user in the database by the given username.
// It returns the matched user and nil error if found, otherwise it returns an empty User model
// and an error indicating the user was not found.
func (u *UserService) SearchUserByUsername(username string) (models.User, error) {

	// search from MongoDB
	filter := bson.M{"username": username}
	found, err := u.Database.Read(filter, &models.User{})
	if err != nil {
		return models.User{}, err
	}

	// found one
	if len(found) > 0 {
		user, ok := found[0].(*models.User)
		if !ok {
			return models.User{}, errors.New("type assertion failed")
		}
		return *user, nil
	}

	// not found
	return models.User{}, errors.New("not found")

}
