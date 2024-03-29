package users

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/Amo-Addai/golang-angular-now-app/golang-fiber-api/src/config"
	"github.com/Amo-Addai/golang-angular-now-app/golang-fiber-api/src/database"
	"github.com/Amo-Addai/golang-angular-now-app/golang-fiber-api/src/helpers"
	"github.com/Amo-Addai/golang-angular-now-app/golang-fiber-api/src/models"
	"github.com/Amo-Addai/golang-angular-now-app/golang-fiber-api/src/services"
	"github.com/Amo-Addai/golang-angular-now-app/golang-fiber-api/src/utils"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

func SignUp(c *fiber.Ctx) error {
	db := database.DB

	data := new(models.User)

	if bodyErr := c.BodyParser(data); bodyErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   bodyErr,
		})
	}

	if ok, err := helpers.ValidateInput(*data); !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err,
		})
	}

	if err := passwordvalidator.Validate(data.Password, 60); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"password": err.Error(),
			},
		})
	}

	if userExists := db.Where(&models.User{Email: data.Email}).First(new(models.User)).RowsAffected; userExists > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Sorry user already exists",
		})
	}

	hash, err := utils.HashPassword(data.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err,
		})
	}

	data.Password = string(hash)

	if createErr := db.Create(&data).Error; createErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"general": "Something went wrong, please try again later. 😕",
		})
	}

	// setting up the email verification
	emailToken := services.GeneralTokens(data.UUID.String(), "verify_email", 24)

	content := services.EmailVerification(utils.EmailVerification{
		Username:   data.FirstName,
		VerifyLink: emailToken,
	})

	helpers.SendEmail(helpers.Payload{
		To:          data.Email,
		Name:        data.FirstName,
		Cc:          "",
		HTMLContent: content,
		Subject:     "Welcome, Please verify your email below",
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Congratulations, your account has been successfully created.",
	})
}

func SignIn(c *fiber.Ctx) error {
	db := database.DB

	data := new(utils.SignIn)

	if bodyErr := c.BodyParser(data); bodyErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   bodyErr,
		})
	}

	if ok, err := helpers.ValidateInput(*data); !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err,
		})
	}

	user := new(models.User)

	if res := db.Where(&models.User{Email: user.Email, Verified: true}).First(user); res.RowsAffected <= 0 {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Incorrect credentials",
		})
	}

	if ok := utils.CheckPasswordHash(data.Password, user.Password); !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Incorrect credentials.",
		})
	}

	// setting up the authorization cookies
	accessToken, refreshToken, accessTime, refreshTime := services.GenerateTokens(user.UUID.String())
	accessCookie, refreshCookie := services.GetAuthCookies(accessToken, refreshToken)
	c.Cookie(accessCookie)
	c.Cookie(refreshCookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"access": fiber.Map{
			"token":  accessToken,
			"expire": accessTime,
		},
		"refresh": fiber.Map{
			"token":  refreshToken,
			"expire": refreshTime,
		},
	})
}

var jwtKey = []byte(config.Config("APP_SECRET"))

func GetAccessToken(c *fiber.Ctx) error {
	db := database.DB

	reToken := new(utils.RefreshToken)

	if bodyErr := c.BodyParser(reToken); bodyErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   bodyErr,
		})
	}

	refreshToken := reToken.RefreshToken

	refreshClaims := new(models.Claims)

	token, _ := jwt.ParseWithClaims(refreshToken, refreshClaims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if res := db.Where(
		"expires_at = ? AND issuer = ? AND audience = ?",
		refreshClaims.ExpiresAt, refreshClaims.Issuer, refreshToken,
	).First(&models.Claims{}); res.RowsAffected <= 0 {
		// no such refresh token exist in the database
		c.ClearCookie("access_token", "refresh_token")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Access Denied",
		})
	}

	if token.Valid {
		if refreshClaims.ExpiresAt < time.Now().Unix() {
			// refresh token is expired
			c.ClearCookie("access_token", "refresh_token")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Access Denied",
			})
		}
	} else {
		// malformed refresh token
		c.ClearCookie("access_token", "refresh_token")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Access Denied",
		})
	}

	if deleteErr := db.Where(
		"issuer = ? AND Audience = ?",
		refreshClaims.Issuer, refreshToken,
	).Delete(refreshClaims).Error; deleteErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Sorry could not delete claims. 😕",
		})
	}

	accessToken, refreshToken, accessTime, refreshTime := services.GenerateTokens(refreshClaims.Issuer)
	accessCookie, refreshCookie := services.GetAuthCookies(accessToken, refreshToken)
	c.Cookie(accessCookie)
	c.Cookie(refreshCookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"access": fiber.Map{
			"token":  accessToken,
			"expire": accessTime,
		},
		"refresh": fiber.Map{
			"token":  refreshToken,
			"expire": refreshTime,
		},
	})
}
