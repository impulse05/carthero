package server

import (
	"carthero/internal/model"

	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Get("/", s.HelloWorldHandler)

	s.App.Get("/health", s.healthHandler)

}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}

func (s *FiberServer) RegisterRiderRoutes() {
	s.App.Get("/riders", s.getRidersHandler)
	s.App.Get("/riders/free", s.getFreeRidersHandler)
	s.App.Post("/riders", s.createRiderHandler)
	s.App.Delete("/riders/:id", s.deleteRiderHandler)
	s.App.Patch("/riders/:id", s.updateRiderStatusHandler)
}

func (s *FiberServer) getRidersHandler(c *fiber.Ctx) error {
	riders, err := s.db.GetRiders()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(riders)
}

func (s *FiberServer) getFreeRidersHandler(c *fiber.Ctx) error {
	riders, err := s.db.GetFreeRiders()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(riders)
}

func (s *FiberServer) createRiderHandler(c *fiber.Ctx) error {
	var rider model.Rider
	if err := c.BodyParser(&rider); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if rider.Name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing name")
	}

	rider, err := s.db.CreateRider(rider)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(rider)
}

func (s *FiberServer) deleteRiderHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = s.db.DeleteRider(id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (s *FiberServer) updateRiderStatusHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var status struct {
		Assigned bool `json:"assigned"`
	}

	if err := c.BodyParser(&status); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = s.db.UpdateRiderStatus(id, status.Assigned)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
