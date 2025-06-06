package db

import (
	"database/sql"
	"fmt"
)

func ProcessImage(db *sql.DB) ([]ImageProcess, error) {
	query := `SELECT images_id, images_name FROM images WHERE process_status_id in ('1','2') ORDER BY images_id ASC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %v", err)
	}
	defer rows.Close()

	var images []ImageProcess

	for rows.Next() {
		var e ImageProcess
		err := rows.Scan(&e.ImageID, &e.ImageName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %v", err)
		}
		images = append(images, e)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return images, nil
}
