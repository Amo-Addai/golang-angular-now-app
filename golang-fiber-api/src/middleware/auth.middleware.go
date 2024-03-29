// SecureAuth returns a middleware which secures all the private routes
package middleware

import (
	"time"

	"github.com/Amo-Addai/golang-angular-now-app/golang-fiber-api/src/config"
	"github.com/Amo-Addai/golang-angular-now-app/golang-fiber-api/src/helpers"
	"github.com/Amo-Addai/golang-angular-now-app/golang-fiber-api/src/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

var jwtKey = []byte(config.Config("APP_SECRET"))

func SecureAuth(c *fiber.Ctx) error {
	accessToken := helpers.ExtractToken(c)

	if len(accessToken) == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Access Denied",
		})
	}

	claims := new(models.Claims)

	token, err := jwt.ParseWithClaims(accessToken, claims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if token.Valid {
		if claims.ExpiresAt < time.Now().Unix() {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"general": "Access Denied",
			})
		}
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			// this is not even a token, we should delete the cookies here
			c.ClearCookie("access_token", "refresh_token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Access Denied",
			})
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Access Denied",
			})
		} else {
			// cannot handle this token
			c.ClearCookie("access_token", "refresh_token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Access Denied",
			})
		}
	}

	c.Locals("id", claims.Issuer)
	return c.Next()
}
