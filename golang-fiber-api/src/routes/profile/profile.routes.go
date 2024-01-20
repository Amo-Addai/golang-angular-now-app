package profile

import (
	"github.com/gofiber/fiber/v2"
	"github.com/Amo-Addai/golang-angular-now-app/golang-fiber-api/src/controllers/users"
)

func ProfileRoutes(router fiber.Router) {
	router.Post("/profile", users.GetProfile)
	router.Put("/update-password", users.UpdatePassword)
}
