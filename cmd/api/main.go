package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/db"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/eventbus"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/firebase"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/glue/deps"
	handlers_users_create "github.com/mrrobotisreal/journey-app-server-refresh/internal/handlers/users/create"
	"log"
	"net/http"
	"os"
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

	fmt.Println("Server up and running on port 6913...")

	if err := http.ListenAndServe(":6913", nil); err != nil {
		log.Fatalf("Failed to start server on port 6913: %v", err)
	}
}
