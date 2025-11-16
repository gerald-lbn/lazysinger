package controllers

import (
	"github.com/gerald-lbn/refrain/pkg/repository"
	"github.com/gerald-lbn/refrain/pkg/services"
	"github.com/gofiber/fiber/v2"
)

type SongsStatController struct {
	container *services.Container
}

func NewSongsStatController(container *services.Container) *SongsStatController {
	return &SongsStatController{
		container: container,
	}
}

// Index returns statistics about songs
func (c *SongsStatController) Index(ctx *fiber.Ctx) error {
	repo := repository.New(c.container.Database)
	stats, err := repo.GetStats(ctx.UserContext())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(stats)
}
