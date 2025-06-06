package router

import (
	"apiPhoto/db"
	"database/sql"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/streadway/amqp"
)

type ImageRequest struct {
	ImageName string `json:"image_name"`
	Actions   string `json:"actions"`
}

func SetupRouteRab(router fiber.Router, conn *sql.DB, mqConn *amqp.Connection) {
	PhotoAPI := router.Group("/")
	{
		PhotoAPI.Post("uploadImage", func(c *fiber.Ctx) error {
			return ImageReq(c, conn, mqConn)
		})
	}
}

func ImageReq(c *fiber.Ctx, conn *sql.DB, mqConn *amqp.Connection) error {
	var req ImageRequest

	// แปลง JSON body เป็น struct
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if req.ImageName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "image_name is required",
		})
	}

	// 🔍 ค้นหาภาพจาก DB ตามชื่อ
	img, err := db.GetImageByName(conn, req.ImageName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if img == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Image not found"})
	}

	// เชื่อมต่อ RabbitMQ
	ch, err := mqConn.Channel()
	if err != nil {
		log.Println("RabbitMQ channel error:", err)
		return c.Status(500).JSON(fiber.Map{"error": "RabbitMQ channel error"})
	}
	defer ch.Close()

	// ประกาศ queue
	q, err := ch.QueueDeclare(
		"image_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Queue declare error"})
	}

	// สร้าง payload
	body, _ := json.Marshal(img)

	// ส่ง message ไป queue
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to send to queue"})
	}

	return c.JSON(fiber.Map{
		"message":    "Image sent to processing queue",
		"images_id":  img.ImageID,
		"image_name": img.ImageName,
	})
}
