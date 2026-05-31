package database

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
    connStr := os.Getenv("DATABASE_URL")  // Render sets this automatically
    if connStr == "" {
        connStr = "host=localhost port=5432 user=postgres password=postgres dbname=blog sslmode=disable"
    }

    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        return err
    }

    _, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS posts (
            id SERIAL PRIMARY KEY,
            title TEXT NOT NULL,
            content TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
    return err
}