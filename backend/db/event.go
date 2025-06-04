package db

import (
	"database/sql"
	"fmt"
)

// func CreateEvent(db *sql.DB, name, description string) error {
func CreateEvent(db *sql.DB, event_name string) error {

	query := `INSERT INTO event (event_name) VALUES (?)`
	_, err := db.Exec(query, event_name)
	if err != nil {
		return fmt.Errorf("failed to insert event: %v", err)
	}
	fmt.Println("Event created successfully")
	return nil
}

type Event struct {
	EventID   int    `json:"event_id"`
	EventName string `json:"event_name"`
	// เพิ่มฟิลด์อื่นๆ ถ้ามี เช่น Description, Location, StartDate ฯลฯ
}

func GetAllEvents(db *sql.DB) ([]Event, error) {
	query := `SELECT event_id, event_name FROM event ORDER BY event_id DESC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %v", err)
	}
	defer rows.Close()

	var events []Event

	for rows.Next() {
		var e Event
		err := rows.Scan(&e.EventID, &e.EventName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %v", err)
		}
		events = append(events, e)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return events, nil
}
