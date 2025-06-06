package main

import (
	"log"
	"time"

	"apiPhoto/db"
	v1 "apiPhoto/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/streadway/amqp"
)

// ฟังก์ชันเชื่อมต่อ RabbitMQ พร้อม retry
func connectRabbitMQ(url string) (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error

	maxRetries := 10                 // จำนวนครั้งสูงสุดที่จะลอง
	retryInterval := 3 * time.Second // เวลาระหว่างลองใหม่

	for i := 0; i < maxRetries; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			log.Println("Connected to RabbitMQ!")
			return conn, nil
		}

		log.Printf("Failed to connect to RabbitMQ (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(retryInterval)
	}

	return nil, err
}

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// เชื่อมต่อฐานข้อมูล
	dbConn := db.ConnectDB()
	defer dbConn.Close()

	// เชื่อมต่อ RabbitMQ พร้อม retry
	mqConn, err := connectRabbitMQ("amqp://skko:skkospiderman@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ after retries: %v", err)
	}
	defer mqConn.Close()

	// ส่ง connection ให้ router หรือ service ที่ต้องใช้
	router := app.Group("/api")
	v1.Setup(router, dbConn, mqConn)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Go langs!")
	})

	log.Fatal(app.Listen(":8000"))
}
