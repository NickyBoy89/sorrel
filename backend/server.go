package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func initDb(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS menus(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		date TIMESTAMP
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
}

func main() {
	db, err := sql.Open("sqlite3", "recipes.db")
	if err != nil {
		log.Fatalf("Error opening db: %v", err)
	}
	defer db.Close()

	// Setup
	initDb(db)

	http.HandleFunc("/api/create-menu", handleCreateMenu)

	log.Print("Serving files...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleCreateMenu(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
