package controllers

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/gerald-lbn/refrain/pkg/repository"
	"github.com/gerald-lbn/refrain/pkg/services"
	dbUtils "github.com/gerald-lbn/refrain/pkg/utils/db"
	"github.com/gofiber/fiber/v2"
)

type SongsController struct {
	container *services.Container
}

func NewSongsController(container *services.Container) *SongsController {
	return &SongsController{
		container: container,
	}
}

// Index returns a list of all songs.
func (c *SongsController) Index(ctx *fiber.Ctx) error {
	repo := repository.New(c.container.Database)
	songs, err := repo.GetAllTracks(ctx.UserContext())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "no songs found",
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.JSON(songs)
}

// Search returns a list of songs matching the given query.
func (c *SongsController) Search(ctx *fiber.Ctx) error {
	query := strings.Trim(ctx.Query("query"), " ")
	if query == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "at least one query parameter must be provided",
		})
	}

	repo := repository.New(c.container.Database)
	songs, err := repo.SearchTracks(ctx.UserContext(), repository.SearchTracksParams{
		Title:  dbUtils.StringToNullString(dbUtils.Like(query)),
		Artist: dbUtils.StringToNullString(dbUtils.Like(query)),
		Album:  dbUtils.StringToNullString(dbUtils.Like(query)),
	})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.JSON(songs)
}

// Show returns a single song by ID.
func (c *SongsController) Show(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid ID",
		})
	}

	repo := repository.New(c.container.Database)
	song, err := repo.GetTrackByID(ctx.UserContext(), intId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "song not found",
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.JSON(song)
}
