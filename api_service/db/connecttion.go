package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

func ConnectDB() *sql.DB {
	// อ่านค่าจาก environment variables หรือใช้ค่า default
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "192.168.0.121"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "3306"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "mysql_121"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "hdcdatarit9esoydld]o8i"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "officedd_photo"
	}

	// MySQL connection string รูปแบบ:
	// user:password@tcp(host:port)/dbname?parseTime=true
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbname)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("cannot open database: %v", err)
	}

	// ทดสอบ connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	fmt.Println("Connected to MySQL database successfully")
	return db
}
