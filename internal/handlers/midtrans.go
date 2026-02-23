package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"midtrans-gateway/internal/midtrans"
	"midtrans-gateway/internal/proxy"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func MidtransWebhook(w http.ResponseWriter, r *http.Request) {
	log.Println("Received Midtrans webhook")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("RAW BODY: %s", body)

	var payload midtrans.WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if !midtrans.VerifySignature(payload) {
		log.Println("Invalid signature")
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}

	r.Body = io.NopCloser(bytes.NewReader(body))
	target := midtrans.ResolveURL(payload.OrderID)
	if err := proxy.Forward(r, target); err != nil {
		log.Println("Unable to forward webhook:", err)
		http.Error(w, "forward failed", http.StatusBadGateway)
		return
	}
	log.Println("forwading to ", target)

	w.WriteHeader(http.StatusOK)
}
