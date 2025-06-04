package router

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/streadway/amqp"
)

type UploadRequest struct {
	ImageIDs []int `json:"image_ids"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func sendToRabbitMQ(imageIDs []int) error {
	// conn, err := amqp.Dial("amqp://skko:skkospiderman@localhost:5672/")
	conn, err := amqp.Dial("amqp://skko:skkospiderman@rabbitmq:5672/")

	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	queueName := "face_jobs"
	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	// ส่งเป็น batch เดียว
	body, err := json.Marshal(map[string]interface{}{
		"image_ids": imageIDs,
	})
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
	return err
}

func SetupImagesAPI(router fiber.Router, conn *sql.DB) {
	PhotoAPI := router.Group("/images")

	PhotoAPI.Post("/notify-upload", func(c *fiber.Ctx) error {
		var req UploadRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON payload",
			})
		}

		if len(req.ImageIDs) == 0 {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "image_id cannot be empty",
			})
		}

		// ตัวอย่างถ้าต้องการ query ข้อมูลรูปจาก DB ก่อน (optional)

		for _, id := range req.ImageIDs {
			var filename string
			err := conn.QueryRow("SELECT images_name FROM images WHERE images_id = ?", id).Scan(&filename)
			if err != nil {
				log.Printf("Image ID %d not found: %v", id, err)
				// อาจจะเลือกข้าม หรือ return error ก็ได้
			} else {
				log.Printf("Image ID %d filename: %s", id, filename)
			}
		}

		// ส่ง message เข้า RabbitMQ
		if err := sendToRabbitMQ(req.ImageIDs); err != nil {
			log.Printf("❌ RabbitMQ error: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to send message to queue",
			})
		}

		return c.JSON(fiber.Map{
			"status": "sent to queue",
		})
	})
}
