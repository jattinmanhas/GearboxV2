package main

import (
    "log"
    router "github.com/jattinmanhas/GearboxV2/services/auth-service/internal/http"
    "net/http"
)

func main() {
    r := router.NewRouter()

    log.Println("Auth Service Running on port :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}