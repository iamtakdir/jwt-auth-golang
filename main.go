package main

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/iamtakdir/jwt-auth-go/database"
	"github.com/iamtakdir/jwt-auth-go/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	connection.Connect()

	app := fiber.New()

	//Cors issue resolved

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	app.Listen("127.0.0.1:3000")
}
