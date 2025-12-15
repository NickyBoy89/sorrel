package main

import (
	"context"
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

	grocerylists "github.com/NickyBoy89/sorrel/backend/grocery_lists"
	"github.com/NickyBoy89/sorrel/backend/internal/auth"
	"github.com/NickyBoy89/sorrel/backend/internal/config"
	"github.com/NickyBoy89/sorrel/backend/internal/db"
	"github.com/NickyBoy89/sorrel/backend/internal/handlers"
	"github.com/NickyBoy89/sorrel/backend/internal/middleware"
	"github.com/SherClockHolmes/webpush-go"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var (
	serverPort = 9031
	debug      = false
	authMethod string
)

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

var conf *config.Config

func init() {
	serveCommand.Flags().IntVar(&serverPort, "port", 9031, "The port to listen on")
	serveCommand.Flags().BoolVar(&debug, "debug", false, "Enable debug logging")
	serveCommand.Flags().StringVar(&authMethod, "auth-method", "keycloak", "The method used to restrict access to non-public API endpoints")
}

var serveCommand = &cobra.Command{
	Use:   "serve",
	Short: "starts the backend process",
	Run: func(cmd *cobra.Command, args []string) {

		if debug {
			log.SetLogLoggerLevel(log.LevelDebug)
			log.Debug("Debug logging enabled!")
		}

		recipesDb, err := sql.Open("sqlite3", "recipes.db")
		if err != nil {
			log.Error("Error opening db", "error", err)
			return
		}
		defer recipesDb.Close()
		db.DB = recipesDb

		// Setup
		if err := initDb(db.DB); err != nil {
			log.Error("error initializing database", "error", err)
			return
		}

		if err := readConfigFile(configFileLocation); err != nil {
			log.Error("error reading config file", "error", err)
			return
		}

		// Initialize auth

		var authMiddleware middleware.AuthMiddleware
		var authHandler handlers.AuthHandler

		switch authMethod {
		case "none":
			log.Warn("Authentification disabled! Please enable this before exposing this service to the internet")

			authHandler = &handlers.NoAuthHandler{}
			authMiddleware = &middleware.NoAuthHandler{}

		case "keycloak":
			log.Info("Tailscale authentification enabled")

			keycloakAuthClient, err := auth.New(context.Background(), conf)
			if err != nil {
				log.Error("error initializing auth client", "error", err)
				return
			}

			authHandler = handlers.NewAuthHandler(keycloakAuthClient)
			authMiddleware = middleware.NewAuthMiddleware(context.Background(), keycloakAuthClient)
		default:
			log.Error("unable to find authentication method", "authMethod", authMethod)
			return
		}

		// Serve frontend
		http.Handle("/", http.FileServer(http.Dir("build/")))

		http.HandleFunc("/api/auth/login", authHandler.LoginHandler)
		http.HandleFunc("/api/auth/callback", authHandler.CallbackHandler)

		// Application
		http.HandleFunc("/api/menu/{menuId}", handleGetMenu)
		http.Handle("/api/menu/{menuId}/edit", authMiddleware.RequireAuth(http.HandlerFunc(handleEditMenu)))
		http.HandleFunc("/api/menu/{menuId}/items", handleGetMenuItems)
		http.Handle("/api/menu/{menuId}/create-item", authMiddleware.RequireAuth(http.HandlerFunc(handleCreateMenuItem)))
		http.HandleFunc("/api/menu/list", handleListMenus)
		http.Handle("/api/menu/create", authMiddleware.RequireAuth(http.HandlerFunc(handleCreateMenu)))
		http.Handle("/api/menu/share", authMiddleware.RequireAuth(http.HandlerFunc(handleShareMenu)))
		http.Handle("/api/items/{itemId}/edit", authMiddleware.RequireAuth(http.HandlerFunc(handleEditMenuItem)))
		http.Handle("/api/items/{itemId}/delete", authMiddleware.RequireAuth(http.HandlerFunc(handleDeleteMenuItem)))

		// Grocery lists
		grocerylists.RegisterHandlers(http.DefaultServeMux)

		// User
		http.HandleFunc("/api/validate-id", handleCheckUserId)
		http.HandleFunc("/api/users", authMiddleware.RequireAuth(http.HandlerFunc(handleListUsers)))
		http.HandleFunc("/api/user", handleGetUser)

		// Push
		http.HandleFunc("/api/push/public-key", handleVAPIDPublicKeyRequest)
		http.HandleFunc("/api/push/subscribe", handlePushSubscription)
		http.HandleFunc("/api/push/validate", handleCheckSubscription)

		log.Info("Serving files...", "port", serverPort)
		log.Error("Error serving data",
			"error", http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil),
		)
	},
}

func readConfigFile(fileName string) error {
	configFile, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	if err != nil {
		// Write default config file if one doesn't exist
		if errors.Is(err, fs.ErrNotExist) {
			log.Info("config file does not exist, creating it now")
			configFile, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				return err
			}

			// Assign to global variables
			priv, pub, err = webpush.GenerateVAPIDKeys()
			if err != nil {
				return err
			}

			conf := config.Config{
				PrivateKey: priv,
				PublicKey:  pub,
			}

			if err := json.NewEncoder(configFile).Encode(conf); err != nil {
				return err
			}

			if _, err := configFile.Seek(0, 0); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	defer configFile.Close()

	// Load values from config

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

	if err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?);", userId).Scan(&isValidId); err != nil {
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
