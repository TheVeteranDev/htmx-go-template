package main

import (
	"log"
	"net/http"

	"github.com/theveterandev/htmx-go-template/database"
	"github.com/theveterandev/htmx-go-template/src/controllers"
	utils "github.com/theveterandev/htmx-go-template/src/utilities"
)

func main() {
	dbi := database.DatabaseInfo{
		Username: utils.GetEnv("DB_USERNAME", "postgres"),
		Password: utils.GetEnv("DB_PASSWORD", "password"),
		Host:     utils.GetEnv("DB_HOST", "localhost"),
		Port:     utils.GetEnv("DB_PORT", "5432"),
		DbName:   utils.GetEnv("DB_NAME", "postgres"),
		SslMode:  utils.GetEnv("DB_SSL_MODE", "disable"),
	}

	db := database.Connect(dbi)
	defer db.Close()

	uc := controllers.InitUserController(db)
	tc := controllers.InitTemplateController()

	http.Handle("/", tc)
	http.Handle("/oauth/v1/", uc)

	host := utils.GetEnv("APP_HOST", "http://localhost")
	port := utils.GetEnv("APP_PORT", ":5000")

	log.Printf("Go server running: %s%s", host, port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
