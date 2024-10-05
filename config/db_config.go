package config

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/adycahyoputro/merchant/utils"
	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	host := utils.GetEnv("DB_HOST")
	port := utils.GetEnv("DB_PORT")
	user := utils.GetEnv("DB_USER")
	password := utils.GetEnv("DB_PASSWORD")
	dbName := utils.GetEnv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("success connect db")
	return db, nil
}

func CloseDB(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Printf("error closing db :%s", err)
	} else {
		log.Println("success db closed")
	}
}