package repository

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Соединение с PostgreSQL
func InitDB() error {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")

	// Строка подключения
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, name)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	// Проверка соединения
	if err := db.Ping(); err != nil {
		return err
	}

	DB = db
	return nil
}
