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

func SendNotificationToUser(userId int, data []byte) error {
	subs, err := db.Query("SELECT id, endpoint, keys_auth, keys_p256dh FROM notification_subscriptions WHERE id = ?", userId)
	if err != nil {
		return err
	}

	for subs.Next() {
		var subscriptionId int
		var sub webpush.Subscription
		if err := subs.Scan(
			&subscriptionId,
			&sub.Endpoint,
			&sub.Keys.Auth,
			&sub.Keys.P256dh,
		); err != nil {
			return err
		}

		resp, err := webpush.SendNotification(data, &sub, &webpush.Options{
			Subscriber:      "example@example.com",
			VAPIDPublicKey:  pub,
			VAPIDPrivateKey: priv,
			TTL:             30,
		})
		if err != nil {
			return err
		}

		// Chrome: https://web.dev/articles/push-notifications-web-push-protocol#response-from-push-service
		switch resp.StatusCode {
		case 201:
		case 429:
			log.Error("error: rate-limited by Push service")
		case 413:
			log.Error("error: payload too large")
		case 400:
			log.Error("error: invalid request to Push service")
		case 410, 404:
			// Not valid, remove
			if _, err := db.Exec("DELETE FROM notification_subscriptions WHERE id = ?", subscriptionId); err != nil {
				return err
			}
		}

		resp.Body.Close()
	}

	return nil
}
