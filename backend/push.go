package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	log "log/slog"
	"net/http"

	"github.com/SherClockHolmes/webpush-go"
)

type PushMessage struct {
	Message   string `json:"data"`
	ActionUrl string `json:"url"`
}

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

	log.Info("incoming subscription", "sub", sub)

	if sub.Sub.Endpoint == "" || sub.Sub.Keys.Auth == "" || sub.Sub.Keys.P256dh == "" {
		http.Error(w, "error: One of the fields was empty", http.StatusBadRequest)
		return
	}

	w.Header().Add("Access-Control-Allow-Origin", "*")

	if _, err := db.Exec("INSERT INTO notification_subscriptions (user_id, endpoint, keys_auth, keys_p256dh) VALUES (?, ?, ?, ?)",
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

func handleCheckSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		http.Error(w, "wrong method", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var sub webpush.Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Debug("Received check for subscription", "sub", sub)

	var hasEndpoint bool

	if err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM notification_subscriptions WHERE endpoint = ?);", sub.Endpoint).Scan(&hasEndpoint); err != nil {
		log.Error("error finding notification endpoint", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !hasEndpoint {
		log.Debug("Subscription does not exist, resubscribing user")
		http.Error(w, "server does not have a subscription with the same endpoint", http.StatusNotFound)
		return
	}

	log.Debug("User was already subscribed")
}

func handleShareMenu(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	type ShareRequest struct {
		MenuId  int   `json:"menuId"`
		UserIds []int `json:"users"`
	}

	var req ShareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var menuName string
	if err := db.QueryRow("SELECT name FROM menus WHERE id = ?", req.MenuId).Scan(&menuName); err != nil {
		log.Error("error fetching menu data", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("A new menu has been shared with you: %s", menuName)

	msg := PushMessage{
		Message:   message,
		ActionUrl: fmt.Sprintf("/menu?menu-id=%d", req.MenuId),
	}

	// Begin a transaction because we're going to read and delete invalid
	// subscriptions at the same time

	tx, err := db.Begin()
	if err != nil {
		log.Error("error beginning transaction", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Prepare a response
	resp := make(map[int]bool)

	for _, userId := range req.UserIds {
		if sent, err := SendNotificationToUser(tx, userId, msg); err != nil {
			log.Error("error sending notification", "error", err)
			http.Error(w, "error sending notification", http.StatusInternalServerError)
			return
		} else {
			resp[userId] = sent
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error("error committing notification changes", "error", err)
		http.Error(w, "error committing notification changes", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("error encoding menu share response", "error", err)
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func SendNotificationToUser(tx *sql.Tx, userId int, message PushMessage) (bool, error) {

	var success bool

	log.Debug("Started sending notifications to user", "userId", userId)

	encodedMessage, err := json.Marshal(message)
	if err != nil {
		return success, err
	}

	subs, err := tx.Query("SELECT id, endpoint, keys_auth, keys_p256dh FROM notification_subscriptions WHERE user_id = ?", userId)
	if err != nil {
		return success, err
	}
	defer subs.Close()

	for subs.Next() {
		log.Info("Reading sub")
		var subscriptionId int
		var sub webpush.Subscription
		if err := subs.Scan(
			&subscriptionId,
			&sub.Endpoint,
			&sub.Keys.Auth,
			&sub.Keys.P256dh,
		); err != nil {
			return success, err
		}

		resp, err := webpush.SendNotification(encodedMessage, &sub, &webpush.Options{
			Subscriber:      "example@example.com",
			VAPIDPublicKey:  pub,
			VAPIDPrivateKey: priv,
			TTL:             30,
		})
		if err != nil {
			return success, err
		}

		log.Debug("Sent notification", "status", resp.StatusCode)

		// Overview: https://pushpad.xyz/blog/list-of-http-status-codes-and-errors-returned-by-web-push-services
		switch resp.StatusCode {
		case 201:
			success = true
		case 429:
			log.Error("error: rate-limited by Push service")
		case 413:
			log.Error("error: payload too large")
		case 400:
			log.Error("error: invalid request to Push service")
		case 410, 404:
			log.Debug("Removed invalid notification")
			// Not valid, remove
			if _, err := tx.Exec("DELETE FROM notification_subscriptions WHERE id = ?", subscriptionId); err != nil {
				return success, err
			}
		}

		resp.Body.Close()
	}

	return success, nil
}
