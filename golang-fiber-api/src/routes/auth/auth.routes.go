package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/Amo-Addai/golang-angular-now-app/golang-fiber-api/src/controllers/users"
)

func AuthRoutes(router fiber.Router) {
	router.Post("/signup", users.SignUp)
	router.Post("/signin", users.SignIn)
	router.Get("/access-token", users.GetAccessToken)
}
