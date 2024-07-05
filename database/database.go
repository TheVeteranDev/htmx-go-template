package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Database struct {
	*sql.DB
}

type DatabaseInfo struct {
	Username string
	Password string
	Host     string
	Port     string
	DbName   string
	SslMode  string
}

func Connect(dbi DatabaseInfo) *Database {
	dsn := "postgres://" + dbi.Username + ":" + dbi.Password + "@" + dbi.Host + ":" + dbi.Port + "/" + dbi.DbName + "?sslmode=" + dbi.SslMode
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to %s:%s/%s database successfully!", dbi.Host, dbi.Port, dbi.DbName)

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Pinged %s:%s/%s database successfully!", dbi.Host, dbi.Port, dbi.DbName)

	return &Database{db}
}
