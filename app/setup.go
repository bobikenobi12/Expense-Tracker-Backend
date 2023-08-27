package app

import (
	"ExpenseTracker/config"
	"ExpenseTracker/database"
	"ExpenseTracker/router"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupAndRunApp() error {
	err := config.LoadENV()
	if err != nil {
		return err
	}

	err = database.StartMongoDb()
	if err != nil {
		return err
	}

	defer database.CloseMongoDb()

	app := fiber.New()

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path} ${latency}\n",
	}))

	router.SetUpRoutes(app)

	config.AddSwaggerRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app.Listen(":" + port)
	log.Println("Server started on port " + port)

	return nil
}
