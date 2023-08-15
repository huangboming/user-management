//go:build wireinject
// +build wireinject

package main

import (
	"usermanagement/internal/handlers"
	"usermanagement/internal/services"

	"github.com/google/wire"
)

func InitializeServer() (*handlers.Server, error) {
	wire.Build(handlers.NewServer, services.NewUserService)
	return &handlers.Server{}, nil
}
