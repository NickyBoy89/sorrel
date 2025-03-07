package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	log "log/slog"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const serverPort = 9031

type Menu struct {
	Id   int       `json:"id"`
	Name string    `json:"name"`
	Date time.Time `json:"date"`
}

type MenuItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func initDb(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS menus(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		date TIMESTAMP
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	query = `CREATE TABLE IF NOT EXISTS items(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		menu_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		description TEXT,
		FOREIGN KEY(menu_id) REFERENCES menus(id)
	);
	`

	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
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

	http.HandleFunc("/api/menu/{menuId}", handleGetMenu)
	http.HandleFunc("/api/menu/{menuId}/items", handleGetMenuItems)
	http.HandleFunc("/api/menu/{menuId}/create-item", handleCreateMenuItem)
	http.HandleFunc("/api/menu/{menuId}/items/edit/{itemId}", handleEditMenuItem)
	http.HandleFunc("/api/menu/list", handleListMenus)
	http.HandleFunc("/api/menu/create", handleCreateMenu)

	log.Info("Serving files...", "port", serverPort)
	log.Error("Error serving data", "error", http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil))
}

func handleGetMenu(w http.ResponseWriter, r *http.Request) {
	menuPathValue := r.PathValue("menuId")
	menuId, err := strconv.Atoi(menuPathValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Access-Control-Allow-Origin", "*")

	var id int
	var name string
	var date time.Time

	if err := db.QueryRow("SELECT id, name, date FROM menus WHERE id = ?", menuId).Scan(&id, &name, &date); err != nil {
		log.Error("error fetching menu", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(Menu{Id: id, Name: name, Date: date}); err != nil {
		log.Info("error encoding menu items to json", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleGetMenuItems(w http.ResponseWriter, r *http.Request) {
	menuItems := []MenuItem{}

	rawMenuId := r.PathValue("menuId")
	menuId, err := strconv.Atoi(rawMenuId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Access-Control-Allow-Origin", "*")

	rows, err := db.Query("SELECT menu_id, name, description FROM items WHERE menu_id = ?", menuId)
	if err != nil {
		log.Error("error fetching menu items", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var menuId int
		var name, description string
		if err := rows.Scan(&menuId, &name, &description); err != nil {
			log.Error("error reading row for menu item", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		menuItems = append(menuItems, MenuItem{Name: name, Description: description})
	}

	if err := json.NewEncoder(w).Encode(menuItems); err != nil {
		log.Error("error encoding menu items", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleCreateMenuItem(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rawMenuId := r.PathValue("menuId")
	menuId, err := strconv.Atoi(rawMenuId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Access-Control-Allow-Origin", "*")

	var name, description string

	for formKey, values := range r.Form {
		var formValue string
		if len(values) > 0 {
			formValue = values[0]
		} else {
			http.Error(w, "error in query string: variable should have at least one value", http.StatusBadRequest)
			return
		}

		switch formKey {
		case "name":
			name = formValue
		case "description":
			description = formValue
		default:
			http.Error(w, "extra value provided: "+formKey, http.StatusBadRequest)
			return
		}
	}

	_, err = db.Exec("INSERT INTO items (menu_id, name, description) VALUES (?, ?, ?)", menuId, name, description)
	if err != nil {
		log.Error("error inserting new menu item", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleEditMenuItem(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rawMenuId := r.PathValue("menuId")
	menuId, err := strconv.Atoi(rawMenuId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_ = menuId

	rawItemId := r.PathValue("menuId")
	itemId, err := strconv.Atoi(rawItemId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var name, description string

	for formKey, values := range r.Form {
		var formValue string
		if len(values) > 0 {
			formValue = values[0]
		} else {
			http.Error(w, "error in query string: variable should have at least one value", http.StatusBadRequest)
			return
		}

		switch formKey {
		case "name":
			name = formValue
		case "description":
			description = formValue
		default:
			http.Error(w, "extra value provided: "+formKey, http.StatusBadRequest)
			return
		}
	}

	w.Header().Add("Access-Control-Allow-Origin", "*")

	if _, err := db.Exec("UPDATE items SET name = ?, description = ? WHERE id = ?", name, description, itemId); err != nil {
		log.Error("error editing menu item", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleCreateMenu(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Access-Control-Allow-Origin", "*")

	var name string
	var date time.Time
	var err error

	for formKey, values := range r.Form {
		var formValue string
		if len(values) > 0 {
			formValue = values[0]
		} else {
			http.Error(w, "error in query string: variable should have at least one value", http.StatusBadRequest)
			return
		}

		switch formKey {
		case "name":
			name = formValue
		case "date":
			date, err = time.Parse(time.DateOnly, formValue)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		default:
			http.Error(w, "extra value provided: "+formKey, http.StatusBadRequest)
			return
		}
	}

	if name == "" || date.IsZero() {
		http.Error(w, "error: missing required options: name, description", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO menus (name, date) VALUES (?, ?)", name, date)
	if err != nil {
		log.Error("Error inserting new menu", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleListMenus(w http.ResponseWriter, r *http.Request) {
	menus := []Menu{}

	rows, err := db.Query("SELECT id, name, date FROM menus")
	if err != nil {
		log.Error("Error fetching menus", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var date time.Time
		if err := rows.Scan(&id, &name, &date); err != nil {
			log.Error("Error reading row", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		menus = append(menus, Menu{
			Id:   id,
			Name: name,
			Date: date,
		})
	}

	if err := json.NewEncoder(w).Encode(menus); err != nil {
		log.Error("Error encoding JSON for menus", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
