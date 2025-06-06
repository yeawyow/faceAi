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

const queueName = "image_queue"

func ImageReq(c *fiber.Ctx, conn *sql.DB, mqConn *amqp.Connection) error {
	var req ImageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	if req.ImageName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "image_name is required"})
	}

	img, err := db.GetImageByName(conn, req.ImageName)
	if err != nil {
		log.Println("DB error:", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if img == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Image not found"})
	}

	ch, err := mqConn.Channel()
	if err != nil {
		log.Println("RabbitMQ channel error:", err)
		return c.Status(500).JSON(fiber.Map{"error": "RabbitMQ channel error"})
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Queue declare error:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Queue declare error"})
	}

	body, err := json.Marshal(img)
	if err != nil {
		log.Println("Marshal error:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to marshal image data"})
	}

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
		log.Println("Failed to publish message:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to send to queue"})
	}

	return c.JSON(fiber.Map{
		"message":    "Image sent to processing queue",
		"images_id":  img.ImageID,
		"image_name": img.ImageName,
	})
}
