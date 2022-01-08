package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	connection "github.com/iamtakdir/jwt-auth-go/database"
	"github.com/iamtakdir/jwt-auth-go/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

//Secret Key

const SecretKey = "secret"

func Register(c *fiber.Ctx) error {

	//Get data from Api
	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return err
	}

	//Encrypt the Password
	pasword, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	if err != nil {
		log.Fatal("error in decrypting password")
	}

	//Grab data and save it in Database

	user := models.User{
		Username: data["username"],
		Email:    data["email"],
		Password: pasword,
	}

	//Connect Database and save

	connection.DB.Create(&user)

	return c.SendStatus(fiber.StatusCreated)

}

//lOGIN FUNCTION

func Login(c *fiber.Ctx) error {

	//Get data from Api
	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return err
	}

	//Find the user from database
	var user models.User

	connection.DB.Where("email= ?", data["email"]).First(&user)
	//Validate User
	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)

		return c.JSON(fiber.Map{
			"message": "User not Found",
		})
	}
	//compare password
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Incorrect password, Try again",
		})
	}

	//Generate JWT token and send it to cookies

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.Email,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 1 day
	})

	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {

		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal Server Error",
		})

	}
	//if everything is ok return authenticated message
	//return c.JSON(fiber.Map{
	//	"message": "Successfully Authenticated",
	//})

	//Store it into cookie

	cookie := fiber.Cookie{
		Name:     "access-token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Success",
	})
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("access-token")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)
	//Querying data with token
	var user models.User
	connection.DB.Where("email = ?", claims.Issuer).First(&user)

	return c.JSON(user)

}

// Logout For logout we need to generate new token
func Logout(c *fiber.Ctx) error {

	cookie := fiber.Cookie{
		Name:     "access-token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "user logout",
	})
}
