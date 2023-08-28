package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"usermanagement/internal/models"

	_ "github.com/go-sql-driver/mysql"

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
	LoginDB()
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

// LoginDB: login database
func (u *UserService) LoginDB() {
	if os.Getenv("MONGO_URI") != "" {
		u.loginMongo()
	} else if os.Getenv("MYSQL_URI") != "" {
		u.loginMySQL()
	} else {
		panic("No database connection")
	}
}

// loginMongo: login MongoDB
func (u *UserService) loginMongo() {
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

// loginMySQL: login MySQL
func (u *UserService) loginMySQL() {
	mysql := models.NewMySQL()
	db, _ := sql.Open("mysql", os.Getenv("MYSQL_URI"))
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MySQL")

	mysql.DB = db
	u.Database = mysql
}

// ----- implement functions for Web API -----

// GetAllUsers get all users in the database and return a list of user
func (u *UserService) GetAllUsers() ([]models.User, error) {

	users := make([]models.User, 0)

	if _, ok := u.Database.(*models.MongoDB); ok {
		// search data from MongoDB
		found, err := u.Database.Read(bson.M{}, func() interface{} { return &models.User{} })
		if err != nil {
			return nil, err
		}
		for _, user := range found {
			u, ok := user.(*models.User)
			if !ok {
				return nil, errors.New("type assertion failed")
			}
			users = append(users, *u)
		}
	} else if _, ok := u.Database.(*models.MySQL); ok {
		// search data from MySQL
		found, err := u.Database.Read("SELECT * FROM users", func() interface{} { return &models.User{} })
		if err != nil {
			return nil, err
		}
		for _, user := range found {
			u, ok := user.(*models.User)
			if !ok {
				return nil, errors.New("type assertion failed")
			}
			users = append(users, *u)
		}
	} else {
		// for unit test
		found, err := u.Database.Read(bson.M{}, func() interface{} { return &models.User{} })
		if err != nil {
			return nil, err
		}
		for _, user := range found {
			u, ok := user.(*models.User)
			if !ok {
				return nil, errors.New("type assertion failed")
			}
			users = append(users, *u)
		}
	}

	return users, nil
}

// CreateUser adds a new user to the UserService.
// If the user with the same username already exists, it returns an error.
func (u *UserService) CreateUser(user models.User) error {

	// if the user already exists, return error
	_, err := u.SearchUserByUsername(user.Username)
	if err == nil {
		return errors.New("user already exist")
	}

	if _, ok := u.Database.(*models.MongoDB); ok {
		// insert to MongoDB
		err = u.Database.Create(user)
		if err != nil {
			return err
		}
	} else if _, ok := u.Database.(*models.MySQL); ok {
		// insert into MySQL
		err = u.Database.Create(fmt.Sprintf("INSERT INTO users VALUES ('%s', '%s', '%s')", user.ID.Hex(), user.Username, user.Password))
		if err != nil {
			return err
		}
	} else {
		// for unit test
		err = u.Database.Create(user)
		if err != nil {
			return err
		}
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

	if _, ok := u.Database.(*models.MongoDB); ok {
		// search from MongoDB
		filter := bson.M{"_id": objectID}
		found, err := u.Database.Read(filter, func() interface{} { return &models.User{} })
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
	} else if _, ok := u.Database.(*models.MySQL); ok {
		found, err := u.Database.Read(fmt.Sprintf("SELECT * FROM users WHERE id = '%s'", ID), func() interface{} { return &models.User{} })
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
	} else {
		// for unit test
		filter := bson.M{"_id": objectID}
		found, err := u.Database.Read(filter, func() interface{} { return &models.User{} })
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
	}

	// not found
	return models.User{}, errors.New("not found")
}

// SearchUserByUsername searches for a user in the database by the given username.
// It returns the matched user and nil error if found, otherwise it returns an empty User model
// and an error indicating the user was not found.
func (u *UserService) SearchUserByUsername(username string) (models.User, error) {

	if _, ok := u.Database.(*models.MongoDB); ok {
		// search from MongoDB
		filter := bson.M{"username": username}
		found, err := u.Database.Read(filter, func() interface{} { return &models.User{} })
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
	} else if _, ok := u.Database.(*models.MySQL); ok {
		found, err := u.Database.Read(fmt.Sprintf("SELECT * FROM users WHERE username = '%s'", username), func() interface{} { return &models.User{} })
		fmt.Println(found)
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
	} else {
		// for unit test

		filter := bson.M{"username": username}
		found, err := u.Database.Read(filter, func() interface{} { return &models.User{} })
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
	}

	// not found
	return models.User{}, errors.New("not found")

}
