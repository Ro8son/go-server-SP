package main

import (
	"database/sql"
	_ "github.com/glebarez/go-sqlite"
	"log"
	"net/http"
	"server/database"
)

type app struct {
	DB    *sql.DB
	CACHE *sql.DB
}

func main() {
	db, err := sql.Open("sqlite", "./database.db")
	if err != nil {
		log.Fatal(err)
	}

	dbCache, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	if err = database.SetupCache(dbCache); err != nil {
		log.Fatal(err)
	}

	app := app{
		DB:    db,
		CACHE: dbCache,
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: app.routes(),
	}

	log.Println("Starting server on port :8080")
	server.ListenAndServe()
}
