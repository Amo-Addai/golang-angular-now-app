package verification

import (
	"github.com/gofiber/fiber/v2"
	"github.com/Amo-Addai/golang-angular-now-app/golang-fiber-api/src/controllers"
)

func VerifyRoute(router fiber.Router) {
	router.Post("/verify/:token", controllers.EmailVerification)
	router.Post("/forgot-password", controllers.ForgotPassword)
	router.Post("/reset-password/:token", controllers.ResetPassword)
}
