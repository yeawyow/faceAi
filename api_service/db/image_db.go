package db

import (
	"database/sql"
)

type Image struct {
	ImageId   int    `json:"images_id"`
	ImageName string `json:"images_name"`
}

func SelectImageByName(db *sql.DB, imageName string) (*Image, error) {
	query := "SELECT images_id, images_name FROM images WHERE images_name = ? LIMIT 1"
	row := db.QueryRow(query, imageName)

	var img Image
	err := row.Scan(&img.ImageId, &img.ImageName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // ไม่มีข้อมูล
		}
		return nil, err
	}

	return &img, nil
}
