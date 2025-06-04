package router

import (
	"apiPhoto/db"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func SetupPhotoAPI(router fiber.Router, conn *sql.DB) {
	PhotoAPI := router.Group("/event")
	{

		PhotoAPI.Post("/create", func(c *fiber.Ctx) error {
			return CreateEventHandler(c, conn)
		})
		PhotoAPI.Get("/fetchall", func(c *fiber.Ctx) error {
			return FetchallEventHandler(c, conn)
		})

	}
}

func CreateEventHandler(c *fiber.Ctx, conn *sql.DB) error {
	eventName := c.FormValue("event_name")
	if eventName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "event_name is required"})
	}

	err := db.CreateEvent(conn, eventName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"eventname": eventName})
}
func FetchallEventHandler(c *fiber.Ctx, conn *sql.DB) error {
	events, err := db.GetAllEvents(conn) // ดึงข้อมูลทั้งหมด
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// ส่งข้อมูล events กลับเป็น JSON
	return c.JSON(fiber.Map{"events": events})
}
