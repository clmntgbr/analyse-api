package main

import (
	"go-api/cmd/api/wire"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
)

func setupRoutes(app *fiber.App, container *wire.Container) {
	setupHealthChecks(app)
	setupWebhooks(app, container)
	setupAPIRoutes(app, container)
}

func setupWebhooks(app *fiber.App, container *wire.Container) {
	webhooks := app.Group("/webhook")

	webhooks.Post("/clerk", container.ClerkMiddleware.Protected(), container.ClerkHandler.Execute)
}

func setupHealthChecks(app *fiber.App) {
	app.Get(healthcheck.LivenessEndpoint, healthcheck.New())
	app.Get(healthcheck.ReadinessEndpoint, healthcheck.New())
	app.Get(healthcheck.StartupEndpoint, healthcheck.New())
}

func setupAPIRoutes(app *fiber.App, container *wire.Container) {
	api := app.Group("/api")

	api.Use(container.AuthenticateMiddleware.Protected())
	setupUsersRoutes(api, container)
}

func setupUsersRoutes(api fiber.Router, container *wire.Container) {
	api.Get("/users/me", container.UserHandler.GetUser)
}
