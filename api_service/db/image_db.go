package db

import (
	"database/sql"
	"fmt"
)

type Imagep struct {
	ImageID   int    `json:"images_id"`
	ImageName string `json:"images_name"`
}

func GetImageByName(db *sql.DB, imageName string) (*Imagep, error) {
	query := `SELECT images_id, images_name FROM events WHERE images_name = ?`

	row := db.QueryRow(query, imageName)

	var img Imagep
	err := row.Scan(&img.ImageID, &img.ImageName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // ไม่พบข้อมูล
		}
		return nil, fmt.Errorf("failed2 to query imagesdfsf: %w", err)
	}

	return &img, nil
}
