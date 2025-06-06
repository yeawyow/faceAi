package db

import (
	"database/sql"
)

type image struct {
	ImageName string `json:"image_naem"`
}

func SelectImageByName(db *sql.DB, imageName string) (*image, error) {
	query := "SELECT  image_name FROM images WHERE image_name = ? LIMIT 1"
	row := db.QueryRow(query, imageName)

	var img image
	err := row.Scan(&img.ImageName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // ไม่มีข้อมูล
		}
		return nil, err
	}

	return &img, nil
}
