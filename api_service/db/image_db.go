package db

import (
	"database/sql"
	"fmt"
)

type Image struct {
	ImageID   int    `json:"images_id"`
	ImageName string `json:"images_name"`
}

func GetImageByName(db *sql.DB, imageName string) (*Image, error) {
	query := `SELECT images_id, images_name FROM images WHERE images_name = ?`

	row := db.QueryRow(query, imageName)

	var img Image
	err := row.Scan(&img.ImageID, &img.ImageName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // ไม่พบข้อมูล
		}
		return nil, fmt.Errorf("failed2 to query image: %w", err)
	}

	return &img, nil
}
