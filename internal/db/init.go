package db

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/journey_app_refresh?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("DBU"), os.Getenv("DBP"))
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("sql.Open error: %w", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("db.Ping error: %w", err)
	}

	if err := createTables(); err != nil {
		fmt.Println("createTables error occurred: ", err)
		return err
	}

	fmt.Println("MySQL connected and schema ensured.")

	Repo = &Repository{DB}

	return nil
}

func createTables() error {
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		user_id BIGINT AUTO_INCREMENT PRIMARY KEY,
		fb_id VARCHAR(255) NOT NULL UNIQUE,
		email VARCHAR(320) NOT NULL UNIQUE,
		username VARCHAR(255),
		api_key CHAR(36) NOT NULL UNIQUE,
		api_key_created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		api_key_last_used_at DATETIME,
		api_key_expires_at DATETIME,
		font VARCHAR(50),
	    theme ENUM('light', 'dark') DEFAULT 'dark',
	    created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_user_fb   (fb_id),
		INDEX idx_user_email(email)
	);`

	entriesTable := `
	CREATE TABLE IF NOT EXISTS entries (
		entry_id VARCHAR(36) NOT NULL PRIMARY KEY,
		user_id BIGINT NOT NULL,
		text TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		last_updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
		INDEX idx_user_created (user_id, created_at)
	);`

	entryLocationsTable := `
	CREATE TABLE IF NOT EXISTS entry_locations (
		location_id BIGINT AUTO_INCREMENT PRIMARY KEY,
		entry_id VARCHAR(36) NOT NULL,
		latitude DOUBLE NOT NULL,
		longitude DOUBLE NOT NULL,
		display_name VARCHAR(255),
		FOREIGN KEY (entry_id) REFERENCES entries(entry_id) ON DELETE CASCADE,
		INDEX idx_entry_id (entry_id)
	);`

	entryTagsTable := `
	CREATE TABLE IF NOT EXISTS entry_tags (
		tag_id BIGINT AUTO_INCREMENT PRIMARY KEY,
		entry_id VARCHAR(36) NOT NULL,
		tag_key VARCHAR(255) NOT NULL,
		tag_value VARCHAR(255),
		FOREIGN KEY (entry_id) REFERENCES entries(entry_id) ON DELETE CASCADE,
		INDEX idx_entry_tag (entry_id, tag_key),
		INDEX idx_tag_key (tag_key)
	);`

	entryImagesTable := `
	CREATE TABLE IF NOT EXISTS entry_images (
		image_id BIGINT AUTO_INCREMENT PRIMARY KEY,
		entry_id VARCHAR(36) NOT NULL,
		image_url VARCHAR(320) NOT NULL,
		FOREIGN KEY (entry_id) REFERENCES entries(entry_id) ON DELETE CASCADE,
		INDEX idx_entry_id (entry_id)
	);`

	// Analytics
	analyticsEventsTable := `
	CREATE TABLE IF NOT EXISTS analytics_events (
		event_id BIGINT AUTO_INCREMENT PRIMARY KEY,
		user_id BIGINT NOT NULL,
		fb_id VARCHAR(255) NOT NULL,
		event_type VARCHAR(100) NOT NULL,
		object_type VARCHAR(100),
		object_id BIGINT,
		event_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		meta_data JSON,
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
		INDEX idx_user_event_time (user_id, event_time),
		INDEX idx_event_type (event_type)
	);`

	// Indexes
	indexQueries := []string{
		`CREATE FULLTEXT INDEX idx_entries_text ON entries(text)`,
		`CREATE INDEX idx_entry_locations_entry_id ON entry_locations(entry_id)`,
		`CREATE INDEX idx_entry_tags_entry_id_key ON entry_tags(entry_id, tag_key)`,
		`CREATE INDEX idx_entry_images_entry_id ON entry_images(entry_id)`,
	}

	if _, err := DB.Exec(usersTable); err != nil {
		return err
	}
	if _, err := DB.Exec(entriesTable); err != nil {
		return err
	}
	if _, err := DB.Exec(entryLocationsTable); err != nil {
		return err
	}
	if _, err := DB.Exec(entryTagsTable); err != nil {
		return err
	}
	if _, err := DB.Exec(entryImagesTable); err != nil {
		return err
	}
	if _, err := DB.Exec(analyticsEventsTable); err != nil {
		return err
	}
	for _, query := range indexQueries {
		_, err := DB.Exec(query)
		if err != nil {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1061 {
				log.Printf("Index already exists, skipping: %s", query)
				continue
			}
			log.Printf("Error executing query: %s, error: %v", query, err)
			return err
		}
	}

	return nil
}
