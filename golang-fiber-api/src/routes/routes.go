package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/Amo-Addai/golang-angular-now-app/golang-fiber-api/src/middleware"
	"github.com/Amo-Addai/golang-angular-now-app/golang-fiber-api/src/routes/auth"
	"github.com/Amo-Addai/golang-angular-now-app/golang-fiber-api/src/routes/profile"
	"github.com/Amo-Addai/golang-angular-now-app/golang-fiber-api/src/routes/verification"
	"github.com/Amo-Addai/golang-angular-now-app/golang-fiber-api/src/routes/welcome"
)

func RouteSetup(app *fiber.App) {
	api := app.Group("/api", logger.New())
	v1 := api.Group("auth")

	protected := api.Group("/private", middleware.SecureAuth)

	welcome.Welcome(v1)
	auth.AuthRoutes(v1)
	verification.VerifyRoute(v1)
	profile.ProfileRoutes(protected)
}
