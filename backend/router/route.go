package router

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func Setup(v1 fiber.Router, conn *sql.DB) {
	SetupImagesAPI(v1, conn)
	SetupFaceResultAPI(v1, conn) // POST /api/face-result ✅ เพิ่มอันนี้

}
