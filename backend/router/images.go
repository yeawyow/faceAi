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

func sendToRabbitMQ(imageIDs []int) error {
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

	body, err := json.Marshal(map[string]interface{}{
		"image_ids": imageIDs,
	})
	if err != nil {
		return err
	}

	log.Printf("📤 Sending to RabbitMQ queue '%s': %s", queueName, string(body))

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

		for _, id := range req.ImageIDs {
			var filename string
			err := conn.QueryRow("SELECT images_name FROM images WHERE images_id = ?", id).Scan(&filename)
			if err != nil {
				log.Printf("Image ID %d not found: %v", id, err)
			} else {
				log.Printf("Image ID %d filename: %s", id, filename)
			}
		}

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

	PhotoAPI.Post("/face-result", FaceResultHandler(conn))
}

func FaceResultHandler(conn *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var result FaceResult

		if err := c.BodyParser(&result); err != nil {
			log.Printf("❌ Invalid JSON from AI: %v", err)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid face result format",
			})
		}

		log.Printf("✅ Received face result for imageId: %d", result.ImageID)
		log.Printf("📦 Face count: %d", result.FaceCount)

		facesJSON, err := json.Marshal(result.Faces)
		if err != nil {
			log.Printf("❌ Failed to marshal faces: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to process faces data",
			})
		}

		_, err = conn.Exec(`
			INSERT INTO face_results (image_id, face_count, faces)
			VALUES (?, ?, ?)
			ON DUPLICATE KEY UPDATE
				face_count = VALUES(face_count),
				faces = VALUES(faces)
		`, result.ImageID, result.FaceCount, facesJSON)

		if err != nil {
			log.Printf("❌ DB insert/update error: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save face result",
			})
		}

		return c.JSON(fiber.Map{
			"status":  "face result saved",
			"imageId": result.ImageID,
		})
	}
}
