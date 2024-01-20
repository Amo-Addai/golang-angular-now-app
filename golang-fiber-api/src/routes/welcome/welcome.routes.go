package welcome

import (
	"github.com/gofiber/fiber/v2"
	"github.com/Amo-Addai/golang-angular-now-app/golang-fiber-api/src/controllers"
)

func Welcome(router fiber.Router) {
	router.Get("/", controllers.Welcome)
}
