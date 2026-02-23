package main

import (
	"log"
	"net/http"
	"os"

	"midtrans-gateway/internal/handlers"
	"midtrans-gateway/internal/midtrans"

	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("ENV") == "production" {
		godotenv.Load(".env.production")
	} else {
		godotenv.Load()
	}
	midtrans.LoadConfig()
	for _, item := range midtrans.WebhookConfig.URLs {
		log.Println("Configured URL:", item.Code, "->", item.URL)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/webhooks/midtrans", handlers.MidtransWebhook)
	mux.HandleFunc("/up", handlers.HealthCheck)

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
