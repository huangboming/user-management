package models

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type CURDInterface interface {
	Create(interface{}) error
	Read(interface{}, func() interface{}) ([]interface{}, error)
	Update(interface{}, interface{}) error
	Delete(interface{}) error
}

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
	if err != nil {
		return err
	}
	return nil
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
