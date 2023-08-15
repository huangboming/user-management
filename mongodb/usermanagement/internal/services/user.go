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
	ctx        context.Context
	collection *mongo.Collection
}

// LoginMongo: login mongodb
func (u *UserService) LoginMongo() {
	client, _ := mongo.Connect(u.ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

	u.collection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("user")
}

func NewUserService() *UserService {
	return &UserService{
		collection: nil,
		ctx:        context.Background(),
	}
}

// GetAllUsers get all users in the database and return a list of user
func (u *UserService) GetAllUsers() ([]models.User, error) {
	// select data from mongodb
	cur, err := u.collection.Find(u.ctx, bson.M{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cur.Close(u.ctx)

	users := make([]models.User, 0)
	for cur.Next(u.ctx) {
		var user models.User
		cur.Decode(&user)
		users = append(users, user)
	}
	return users, nil
}

// CreateUser adds a new user to the UserService.
// If the user with the same username already exists, it returns an error.
func (u *UserService) CreateUser(user models.User) error {

	// if the user already exists, return error
	_, err := u.SearchUserByUsername(user.Username)
	if err == nil {
		return errors.New("user already exsit")
	}

	_, err = u.collection.InsertOne(u.ctx, user)
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
	cur, err := u.collection.Find(u.ctx, filter)
	if err != nil {
		return models.User{}, err
	}

	var user models.User
	if cur.Next(u.ctx) { // found one
		err = cur.Decode(&user)
		if err != nil {
			return models.User{}, err
		}
		return user, nil
	} else { // not found
		return models.User{}, errors.New("not found")
	}
}

// SearchUserByUsername searches for a user in the database by the given username.
// It returns the matched user and nil error if found, otherwise it returns an empty User model
// and an error indicating the user was not found.
func (u *UserService) SearchUserByUsername(username string) (models.User, error) {

	// search from MongoDB
	filter := bson.M{"username": username}
	cur, err := u.collection.Find(u.ctx, filter)
	if err != nil {
		return models.User{}, err
	}

	var user models.User
	if cur.Next(u.ctx) { // found one
		err = cur.Decode(&user)
		if err != nil {
			return models.User{}, err
		}
		return user, nil
	} else { // not found
		return models.User{}, errors.New("not found")
	}

}
