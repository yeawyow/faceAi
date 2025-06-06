package router

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/streadway/amqp"
)

func Setup(v1 fiber.Router, conn *sql.DB, mqConn *amqp.Connection) {
	SetupImagesAPI(v1, conn)
	SetupPhotoAPI(v1, conn)
	SetupRouteRab(v1, conn, mqConn)
}
