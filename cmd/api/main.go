package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/db"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/eventbus"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/firebase"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/glue/deps"
	handlers_entries_create "github.com/mrrobotisreal/journey-app-server-refresh/internal/handlers/entries/create"
	handlers_entries_read "github.com/mrrobotisreal/journey-app-server-refresh/internal/handlers/entries/read"
	handlers_users_create "github.com/mrrobotisreal/journey-app-server-refresh/internal/handlers/users/create"
	handlers_users_login "github.com/mrrobotisreal/journey-app-server-refresh/internal/handlers/users/login"
)

func main() {
	_ = godotenv.Load()

	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to init MySQL: %v", err)
	}
	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {
			log.Fatalf("Failed to close MySQL database: %v", err)
		}
	}(db.DB)

	firebase.InitFB()

	brokers := []string{os.Getenv("KAFKA_BROKERS")}
	bus := eventbus.NewBus(brokers)

	deps.Bus = bus

	http.HandleFunc("/users/create", handlers_users_create.CreateAccount)
	http.HandleFunc("/users/login", handlers_users_login.Login)
	http.HandleFunc("/entries/create", handlers_entries_create.CreateEntry)
	http.HandleFunc("/entries/list", handlers_entries_read.ListEntries)
	fmt.Println("Server up and running on port 6913...")

	if err := http.ListenAndServe(":6913", nil); err != nil {
		log.Fatalf("Failed to start server on port 6913: %v", err)
	}
}
