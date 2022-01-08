package controllers

import (
	"github.com/gofiber/fiber/v2"
	connection "github.com/iamtakdir/jwt-auth-go/database"
	"github.com/iamtakdir/jwt-auth-go/models"
	"golang.org/x/crypto/bcrypt"
	"log"
)

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

		return c.JSONP(fiber.Map{
			"message": "User not Found",
		})
	}
	//compare password
	bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"]))

	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Incorrect password, Try again",
		})
	}

	return c.JSON(&user)
}
