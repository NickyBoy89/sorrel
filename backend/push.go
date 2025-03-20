package main

import (
	"encoding/json"
	log "log/slog"
	"net/http"

	"github.com/SherClockHolmes/webpush-go"
)

var (
	priv = ""
	pub  = ""
)

func handleVAPIDPublicKeyRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	w.Write([]byte(pub))
}

func handlePushSubscription(w http.ResponseWriter, r *http.Request) {

	var sub struct {
		UserID int                  `json:"userId"`
		Sub    webpush.Subscription `json:"sub"`
	}

	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if sub.Sub.Endpoint == "" || sub.Sub.Keys.Auth == "" || sub.Sub.Keys.P256dh == "" {
		http.Error(w, "error: One of the fields was empty", http.StatusBadRequest)
		return
	}

	w.Header().Add("Access-Control-Allow-Origin", "*")

	if _, err := db.Exec("INSERT INTO notification_subscriptions (user_id, endpoint, keys_auto, keys_p256dh) VALUES (?, ?, ?, ?)",
		sub.UserID,
		sub.Sub.Endpoint,
		sub.Sub.Keys.Auth,
		sub.Sub.Keys.P256dh,
	); err != nil {
		log.Error("error adding subcription", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("received new subscription", "userId", sub.UserID)
}
