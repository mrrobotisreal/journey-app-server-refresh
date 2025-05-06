package main

import (
	"database/sql"
	"fmt"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/db"
	handlers_users_create "github.com/mrrobotisreal/journey-app-server-refresh/internal/handlers/users/create"
	"log"
	"net/http"
)

func main() {
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to init MySQL: %v", err)
	}
	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {
			log.Fatalf("Failed to close MySQL database: %v", err)
		}
	}(db.DB)

	http.HandleFunc("/users/create", handlers_users_create.CreateAccount)

	fmt.Println("Server up and running on port 6913...")

	if err := http.ListenAndServe(":6913", nil); err != nil {
		log.Fatalf("Failed to start server on port 6913: %v", err)
	}
}
