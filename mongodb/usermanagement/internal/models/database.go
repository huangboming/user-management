package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CURDInterface interface {
	Create(interface{}) error
	Read(interface{}, func() interface{}) ([]interface{}, error)
	Update(interface{}, interface{}) error
	Delete(interface{}) error
}

// ----- MongoDB -----

type MongoDB struct {
	Ctx        context.Context
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewMongoDB() *MongoDB {
	return &MongoDB{
		Ctx:        context.Background(),
		Client:     nil,
		Collection: nil,
	}
}

func (m *MongoDB) Create(item interface{}) error {
	_, err := m.Collection.InsertOne(m.Ctx, item)
	return err
}

func (m *MongoDB) Read(filter interface{}, callback func() interface{}) ([]interface{}, error) {
	// callback is a function that returns an empty interface
	// this is used to create a new instance of the struct that we want to decode the result into
	// e.g. callback := func() interface{} { return &models.User{} }

	cur, err := m.Collection.Find(m.Ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(m.Ctx)

	items := make([]interface{}, 0)
	for cur.Next(m.Ctx) {
		result := callback()
		err := cur.Decode(result)
		if err != nil {
			return nil, err
		}
		items = append(items, result)
	}
	fmt.Println(items)
	return items, nil
}

func (m *MongoDB) Update(filter interface{}, update interface{}) error {
	_, err := m.Collection.UpdateMany(m.Ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) Delete(filter interface{}) error {
	_, err := m.Collection.DeleteMany(m.Ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

// ----- MySQL -----

type MySQL struct {
	DB *sql.DB
}

func NewMySQL() *MySQL {
	return &MySQL{
		DB: nil,
	}
}

func (m *MySQL) Create(query interface{}) error {
	q, ok := query.(string)
	if !ok {
		return errors.New("type assertion failed")
	}
	_, err := m.DB.Exec(q)
	if err != nil {
		return err
	}
	return nil
}

func (m *MySQL) Read(query interface{}, callback func() interface{}) ([]interface{}, error) {
	s, ok := query.(string)

	if !ok {
		return nil, errors.New("type assertion failed")
	}

	rows, err := m.DB.Query(s)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]interface{}, 0)
	for rows.Next() {
		result := callback()
		switch v := result.(type) {

		case *User:
			var idString string
			err := rows.Scan(&idString, &v.Username, &v.Password)
			if err != nil {
				return nil, err
			}

			v.ID, err = primitive.ObjectIDFromHex(idString)
			if err != nil {
				return nil, err
			}

			items = append(items, v)

		default:
			return nil, errors.New("unknown type")

		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// 这里因为API没有用到Update和Delete，所以没有实现
// 下面的代码都没有用，如果以后有API需要，再实现
func (m *MySQL) Update(filter interface{}, update interface{}) error {
	user, ok := filter.(User)
	if !ok {
		return errors.New("type assertion failed")
	}
	_, err := m.DB.Exec("UPDATE user SET password = ? WHERE username = ?", user.Password, user.Username)
	if err != nil {
		return err
	}
	return nil
}

func (m *MySQL) Delete(filter interface{}) error {
	user, ok := filter.(User)
	if !ok {
		return errors.New("type assertion failed")
	}
	_, err := m.DB.Exec("DELETE FROM user WHERE username = ?", user.Username)
	if err != nil {
		return err
	}
	return nil
}
