package repository

import (
	"github.com/Asylann/gRPC_Demo/server/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"time"
)

var db *sqlx.DB

func InitDB(config config.Config) *sqlx.DB {
	var err error
	db, err = sqlx.Open("postgres", config.DatabaseConnection)
	if err != nil {
		log.Fatal(err.Error())
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err.Error())
	}

	/*RunMigration(db)*/

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(4 * time.Minute)

	log.Println("Postgres DB is connected!")
	return db
}
