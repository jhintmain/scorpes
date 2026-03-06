package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hooneun/scorpes/internal/api"
)

func main() {

	r := api.NewRouter()

	r.Use(api.LoggerMiddleware)

	r.GET("/health", helloHandler)

	if err := http.ListenAndServe(":8090", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Health check successful"})
}
