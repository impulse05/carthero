package server

import (
	"github.com/gofiber/fiber/v2"

	"carthero/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "carthero",
			AppName:      "carthero",
		}),

		db: database.New(),
	}

	return server
}
