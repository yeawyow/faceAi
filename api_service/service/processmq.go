package service

import (
	"apiPhoto/db"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type Image struct {
	ImageID   int    `json:"image_id"`
	ImageName string `json:"image_name"`
}

type ImagePayload struct {
	Images []Image `json:"images"`
}

func SendProcessMq(images []db.ImageProcess) error {
	if len(images) == 0 {
		log.Println("⚠️ No images to process")
		return nil
	}

	var mappedImages []Image
	for _, img := range images {
		mappedImages = append(mappedImages, Image{
			ImageID:   img.ImageID,
			ImageName: img.ImageName,
		})
	}

	payload := ImagePayload{
		Images: mappedImages,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal failed: %w", err)
	}

	log.Printf("📦 JSON Payload: %s\n", string(body))

	conn, err := amqp.Dial("amqp://skko:skkospiderman@rabbitmq:5672/")
	if err != nil {
		return fmt.Errorf("rabbitmq dial failed: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("⚠️ Error closing connection: %v\n", err)
		}
	}()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("channel failed: %w", err)
	}
	defer func() {
		if err := ch.Close(); err != nil {
			log.Printf("⚠️ Error closing channel: %v\n", err)
		}
	}()

	q, err := ch.QueueDeclare(
		"face_images_queue", // queue name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return fmt.Errorf("queue declare failed: %w", err)
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key (queue name)
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	log.Println("✅ ส่ง images ไปยัง RabbitMQ แล้ว")
	return nil
}
