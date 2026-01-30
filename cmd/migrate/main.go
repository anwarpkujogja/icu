package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/hamilton/icu-app/pkg/util/db"
	"github.com/joho/godotenv"
)

func main() {
	// flags
	up := flag.Bool("up", false, "Migrate Up")
	down := flag.Bool("down", false, "Migrate Down (1 step)")
	force := flag.Int("force", 0, "Force version")
	flag.Parse()

	// 1. Load Env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, checking system env")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	// 2. Setup Migration
	// Resolve absolute path for etc/migrations to avoid file:// relative path issues
	cwd, _ := os.Getwd()
	migrationPath := filepath.Join(cwd, "etc", "migrations")
	
	// Check if path exists
	if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
		log.Printf("Migration path not found at: %s. Trying relative...", migrationPath)
		migrationPath = "etc/migrations"
	}
	
	// Create file:// URL. On Windows this needs to be handled carefully, 
	// but golang-migrate usually accepts "file://C:/..." or "file://./etc/migrations"
	// Let's use simple string concat with forward slashes for safety in URL
	migrationSource := "file://" + filepath.ToSlash(migrationPath)
	log.Println("Migration Source:", migrationSource)

	url := db.ConstructMigrationUrl(dbUser, dbPass, dbHost, dbPort, dbName)
	
	m, err := db.NewMigration(migrationSource, url)
	if err != nil {
		log.Fatalf("Failed to init migration: %v", err)
	}

	// 3. Execute
	if *up {
		if err := m.Up(); err != nil {
			log.Fatalf("Migration Up failed: %v", err)
		}
	} else if *down {
		if err := m.MigrateOneStepDown(); err != nil {
			log.Fatalf("Migration Down failed: %v", err)
		}
	} else if *force != 0 {
		if err := m.ForceMigrate(*force); err != nil {
			log.Fatalf("Force migration failed: %v", err)
		}
	} else {
		log.Println("No action specified. Use -up, -down, or -force.")
	}
}
