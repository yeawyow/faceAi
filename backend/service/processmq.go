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

	conn, err := amqp.Dial("amqp://skko:skkospiderman@rabbitmq:5672/")
	if err != nil {
		return fmt.Errorf("rabbitmq dial failed: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("channel failed: %w", err)
	}

	q, err := ch.QueueDeclare(
		"face_images_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue declare failed: %w", err)
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	log.Println("✅ ส่ง images ไปยัง RabbitMQ แล้ว")
	defer ch.Close()
	return nil

}
