package router

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type FaceResult struct {
	ImageID   int           `json:"imageId"`
	FaceCount int           `json:"face_count"`
	Faces     []FaceDetails `json:"faces"`
}

type FaceDetails struct {
	BBox      []float64   `json:"bbox"`      // [x1, y1, x2, y2]
	Embedding []float64   `json:"embedding"` // 128D หรือ 512D
	Landmark  [][]float64 `json:"landmark"`  // 5 หรือ 106 keypoints
}

func SetupFaceResultAPI(router fiber.Router, conn *sql.DB) {
	router.Post("/face-result", func(c *fiber.Ctx) error {
		var result FaceResult

		if err := c.BodyParser(&result); err != nil {
			log.Printf("❌ ไม่สามารถแปลง JSON: %v", err)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON format",
			})
		}

		log.Printf("📬 ได้รับผลลัพธ์จาก AI สำหรับ imageId=%d (จำนวนใบหน้า: %d)", result.ImageID, result.FaceCount)

		// คุณสามารถบันทึก result ลงฐานข้อมูล หรือทำอย่างอื่นต่อได้
		// เช่น แสดง embedding/bbox:
		for i, face := range result.Faces {
			log.Printf("🔹 Face #%d: BBox: %v", i+1, face.BBox)
			log.Printf("🔹 Landmark: %v", face.Landmark)
			log.Printf("🔹 Embedding length: %d", len(face.Embedding))
		}

		// ส่ง response กลับ
		return c.JSON(fiber.Map{
			"status":  "received",
			"imageId": result.ImageID,
			"faces":   result.FaceCount,
		})
	})
}
