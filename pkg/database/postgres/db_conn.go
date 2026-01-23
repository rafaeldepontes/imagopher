package postgres

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func GetDB() *sql.DB {
	if db != nil {
		return db
	}
	_ = initDB()
	return db
}

func initDB() error {
	database, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println("[ERROR] Could not open the connect to the database:", err)
		return err
	}

	if err := db.Ping(); err != nil {
		log.Println("[ERROR] Could not verify the connection:", err)
		return err
	}

	db = database
	return nil
}
