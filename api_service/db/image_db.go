package db

import (
	"database/sql"
)

type Image struct {
	ImageName string `json:"image_naem"`
}

func SelectImageByName(db *sql.DB, imageName string) (*Image, error) {
	query := "SELECT id, image_name, url FROM images WHERE image_name = ? LIMIT 1"
	row := db.QueryRow(query, imageName)

	var img Image
	err := row.Scan(&img.ImageName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // ไม่มีข้อมูล
		}
		return nil, err
	}

	return &img, nil
}
