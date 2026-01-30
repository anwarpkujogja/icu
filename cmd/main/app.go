package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hamilton/icu-app/pkg/icu/delivery/http"
	"github.com/hamilton/icu-app/pkg/icu/repository"
	"github.com/hamilton/icu-app/pkg/icu/usecase"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	// 1. Load Env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)
	log.Printf("DEBUG: Connecting to DB Host: %s, DSN (masked): host=%s dbname=%s", dbHost, dbHost, dbName)

	// 2. Connect DB
	dbConn, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// 3. Setup Architecture Layers
	timeoutContext := time.Duration(10) * time.Second

	repo := repository.NewPostgresRepository(dbConn)
	uc := usecase.NewICUUseCase(repo, timeoutContext)

	// 4. Setup Echo
	e := echo.New()
	http.NewICUHandler(e, uc)

	// 5. Start Server
	log.Printf("Starting server on port %s", appPort)
	log.Fatal(e.Start(":" + appPort))
}
