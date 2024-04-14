package main

import (
	"log"
	"my_app/internal/cache"
	"my_app/internal/db"
	"my_app/internal/server"
	"net/http"
)

func main() {
	log.Printf("Server started")

	db.InitDB()
	defer db.CloseDB()

	cache.InitCache()
	defer cache.CloseCache()

	router := server.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
