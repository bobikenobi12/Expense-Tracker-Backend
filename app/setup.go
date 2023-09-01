package app

import (
	"ExpenseTracker/config"
	"ExpenseTracker/database"
	"ExpenseTracker/router"
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupAndRunApp() error {
	ctx := context.Background()

	if err := config.LoadENV(); err != nil {
		return err
	}

	if err := database.NewDbConn(); err != nil {
		return err
	}

	defer database.CloseConn()

	if err := database.CreateSchema(ctx); err != nil {
		return err
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: config.GlobalErrorHandler,
	})

	config.InitValidations()

	app.Use(recover.New())

	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			env := os.Getenv("GO_ENV")
			return env == "development" || env == ""
		},
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

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
