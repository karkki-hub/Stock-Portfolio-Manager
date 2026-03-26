package database

import (
	"database/sql"
	"fmt"
	"log"

	"karkki-hub/Stock-Portfolio-Manager/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySQL(cfg *config.Config) *sql.DB {

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
