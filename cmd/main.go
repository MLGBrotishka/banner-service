package main

import (
	"log"
	"my_app/internal/db"
	"my_app/internal/server"
	"net/http"
)

func main() {
	log.Printf("Server started")
	db.InitDB()
	router := server.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
