package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	log "log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	_ "github.com/mattn/go-sqlite3"
)

const serverPort = 9031
const configFileLocation = "sorrel-config.json"

type Config struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

type Menu struct {
	Id   int       `json:"id"`
	Name string    `json:"name"`
	Date time.Time `json:"date"`
}

type MenuItem struct {
	Id          int     `json:"id,omitempty"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Section     *string `json:"section,omitempty"`
}

var db *sql.DB

func main() {
	recipesDb, err := sql.Open("sqlite3", "recipes.db")
	if err != nil {
		log.Error("Error opening db", "error", err)
		return
	}
	defer recipesDb.Close()
	db = recipesDb

	// Setup
	if err := initDb(db); err != nil {
		log.Error("error initializing database", "error", err)
		return
	}

	if err := readConfigFile(configFileLocation); err != nil {
		log.Error("error reading config file", "error", err)
		return
	}

	http.Handle("/", http.FileServer(http.Dir("../build")))

	// Application
	http.HandleFunc("/api/menu/{menuId}", handleGetMenu)
	http.HandleFunc("/api/menu/{menuId}/edit", handleEditMenu)
	http.HandleFunc("/api/menu/{menuId}/items", handleGetMenuItems)
	http.HandleFunc("/api/menu/{menuId}/create-item", handleCreateMenuItem)
	http.HandleFunc("/api/menu/list", handleListMenus)
	http.HandleFunc("/api/menu/create", handleCreateMenu)
	http.HandleFunc("/api/menu/share", handleShareMenu)
	http.HandleFunc("/api/items/{itemId}/edit", handleEditMenuItem)
	http.HandleFunc("/api/items/{itemId}/delete", handleDeleteMenuItem)

	http.HandleFunc("/api/validate-id", handleCheckUserId)

	// Push
	http.HandleFunc("/api/push/public-key", handleVAPIDPublicKeyRequest)
	http.HandleFunc("/api/push/subscribe", handlePushSubscription)

	log.Info("Serving files...", "port", serverPort)
	log.Error("Error serving data",
		"error", http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil),
	)
}

func readConfigFile(fileName string) error {
	configFile, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			log.Info("config file does not exist, creating it now")
			configFile, err = os.Create(configFileLocation)
			if err != nil {
				return err
			}

			// Assign to global variables
			priv, pub, err = webpush.GenerateVAPIDKeys()
			if err != nil {
				return err
			}

			conf := Config{
				PrivateKey: priv,
				PublicKey:  pub,
			}

			if err := json.NewEncoder(configFile).Encode(conf); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	defer configFile.Close()

	var conf Config
	if err := json.NewDecoder(configFile).Decode(&conf); err != nil {
		return err
	}

	priv = conf.PrivateKey
	pub = conf.PublicKey

	return nil
}

func handleCheckUserId(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId, err := strconv.Atoi(r.Form.Get("userId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Access-Control-Allow-Origin", "*")

	var isValidId bool

	if err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?);", userId).Scan(&isValidId); err != nil {
		log.Error("error testing for user id", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !isValidId {
		http.Error(w, "Unknown or invalid user id", http.StatusUnauthorized)
		return
	}
	r.Body.Close()
}
