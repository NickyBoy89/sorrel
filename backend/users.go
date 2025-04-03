package main

import (
	"encoding/json"
	log "log/slog"
	"net/http"
	"strconv"
)

type User struct {
	Id          int    `json:"id"`
	DisplayName string `json:"display_name"`
}

func handleGetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodGet {
		http.Error(w, "wrong method", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId, err := strconv.Atoi(r.Form.Get("userId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var requestedUser User

	if err := db.QueryRow("SELECT display_name FROM users WHERE id = ?", userId).Scan(&requestedUser.DisplayName); err != nil {
		log.Error("error finding user", "error", err)
		http.Error(w, "error getting user data", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(requestedUser); err != nil {
		log.Error("error encoding user", "error", err)
		http.Error(w, "error encoding user data", http.StatusInternalServerError)
		return
	}
}

func handleListUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	users, err := db.Query("SELECT id, display_name FROM users")
	if err != nil {
		log.Error("error reading users", "error", err)
		http.Error(w, "error listing users", http.StatusInternalServerError)
		return
	}
	defer users.Close()

	var resultUsers []User

	for users.Next() {
		var cur User

		if err := users.Scan(&cur.Id, &cur.DisplayName); err != nil {
			log.Error("error reading user", "error", err)
			http.Error(w, "error reading users", http.StatusInternalServerError)
			return
		}

		resultUsers = append(resultUsers, cur)
	}

	if err := json.NewEncoder(w).Encode(resultUsers); err != nil {
		log.Error("error encoding users", "error", err)
		http.Error(w, "error encoding users", http.StatusInternalServerError)
		return
	}
}
