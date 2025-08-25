package main

import (
    "log"
    router "github.com/jattinmanhas/GearboxV2/services/auth-service/internal/http"
    "net/http"
    "os"

    "github.com/joho/godotenv"
    "github.com/jattinmanhas/GearboxV2/services/auth-service/db"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    // Connect to DB
    database, err := db.Connect()
    if err != nil {
        log.Fatalf("Failed to connect to the database: %v", err)
    }
    defer database.Close()

    r := router.NewRouter(database)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Println("Auth Service Running on port:", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}