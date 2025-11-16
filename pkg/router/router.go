package router

import (
	"time"

	"github.com/gerald-lbn/refrain/pkg/controllers"
	"github.com/gerald-lbn/refrain/pkg/services"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/utils"
)

func BuildRouter(c *services.Container) error {
	// Middlewares
	c.Web.Use(cors.New())
	c.Web.Use(csrf.New(csrf.Config{
		KeyLookup:      "header:X-Csrf-Token",
		CookieName:     "csrf_",
		CookieSameSite: "Lax",
		Expiration:     1 * time.Hour,
		KeyGenerator:   utils.UUIDv4,
	}))
	c.Web.Use(healthcheck.New(healthcheck.Config{
		LivenessEndpoint:  "/live",
		ReadinessEndpoint: "/ready",
	}))
	c.Web.Use(recover.New())

	// API Routers
	c.Web.Get("/api/stats", controllers.NewSongsStatController(c).Index)
	c.Web.Get("/api/tracks", controllers.NewSongsController(c).Index)
	c.Web.Get("/api/tracks/:id", controllers.NewSongsController(c).Show)
	c.Web.Get("/api/search/tracks", controllers.NewSongsController(c).Search)

	return nil
}
