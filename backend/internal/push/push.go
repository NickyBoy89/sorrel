package push

import (
	"database/sql"
	"encoding/json"

	log "log/slog"

	"github.com/SherClockHolmes/webpush-go"
)

var (
	priv = ""
	Pub  = ""
)

func InitializePrivateKey(value string) {
	priv = value
}

type PushMessage struct {
	Message   string `json:"data"`
	ActionUrl string `json:"url"`
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
			VAPIDPublicKey:  Pub,
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
