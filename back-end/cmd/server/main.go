package main

import (
	"log"
	"net/http"

	"go-do-list/internal/config"
	"go-do-list/internal/database"
	"go-do-list/internal/routes"
)

func main() {
	config.LoadEnv()

	db := database.Connect()
	defer db.Close()

	router := routes.SetupRoutes(db)

	log.Println("Server running on :3333")
	log.Fatal(http.ListenAndServe(":3333", router))
}
