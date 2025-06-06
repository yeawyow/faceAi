package router

import (
	"apiPhoto/db"
	"apiPhoto/service"
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
		PhotoAPI.Post("/process", func(c *fiber.Ctx) error {
			return ProcessEventHandler(c, conn)
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

func ProcessEventHandler(c *fiber.Ctx, conn *sql.DB) error {
	action := c.FormValue("action")
	if action == "process" {
		images, err := db.ProcessImage(conn)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		if len(images) == 0 {
			return c.JSON(fiber.Map{"status": "ไม่มีภาพที่ต้องประมวลผล"})
		}
		if err := service.SendProcessMq(images); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"status": "ส่งข้อมูลสำเร็จ", "images": images})
	}
	return c.SendString("No action performed")
}
