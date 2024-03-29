package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/Amo-Addai/golang-angular-now-app/golang-fiber-api/src/config"
	"github.com/Amo-Addai/golang-angular-now-app/golang-fiber-api/src/database"
	"github.com/Amo-Addai/golang-angular-now-app/golang-fiber-api/src/routes"
)

func main() {
	fibreConfig := fiber.Config{
		ServerHeader: config.Config("APP_NAME"),
	}
	app := fiber.New(fibreConfig)

	// default routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "All good here!.. 🚀 ",
		})
	})

	// app config
	app.Use(cors.New())

	port := config.Config("PORT")

	// postgres, mongodb
	database.DBConnection("postgres")

	// routes
	routes.RouteSetup(app)

	conErr := app.Listen(":" + port)

	// connection error
	if conErr != nil {
		panic(conErr)
	}
}
