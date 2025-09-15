package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(128) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX scheduler_date ON scheduler (date);
`

func GetDBAddress() (string, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Note: .env file not found, using system environment variables")
		return "", err
	}
	dbAddress := os.Getenv("TODO_DBFILE")
	if dbAddress == "" {
		dbAddress = "scheduler.db"
		log.Printf("TODO_DBFILE not set, using default: %s", dbAddress)
	}
	return dbAddress, nil
}

func Init(dbFile string) error {

	install := false
	_, err := os.Stat(dbFile)
	if os.IsNotExist(err) {
		install = true
	} else if err != nil {
		return fmt.Errorf("error checking database file: %v", err)
	}

	DB, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	if err := DB.Ping(); err != nil {
		DB.Close()
		DB = nil
		return fmt.Errorf("error connecting to database: %v", err)
	}

	if install {
		_, err = DB.Exec(schema)
		if err != nil {
			DB.Close()
			DB = nil
			os.Remove(dbFile)
			return fmt.Errorf("error creating database schema: %v", err)
		}
		fmt.Printf("Database created and initialized: %s\n", dbFile)
	}

	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
