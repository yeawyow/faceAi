package main

import (
	"database/sql"
	"log"

	"apiPhoto/db"
	v1 "apiPhoto/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var dbConn *sql.DB

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	dbConn = db.ConnectDB()

	defer dbConn.Close()

	// ส่ง conn ให้ router หรือ service ที่ต้องใช้
	router := app.Group("/api")
	v1.Setup(router, dbConn)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Go langs!")
	})

	log.Fatal(app.Listen(":8000"))
}
