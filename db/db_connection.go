package db

import (
	"database/sql"
	"fmt"
	"log"

	"example.com/golang_twitter/config"
	_ "github.com/lib/pq"
)

var dbConn *sql.DB

// DB接続を行う
func ConnectDB(cfg config.Config) {
	dbinfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	var err error
	dbConn, err = sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}

	if err = dbConn.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("DB Connected!")
}

func DbConn() *sql.DB {
	return dbConn
}
