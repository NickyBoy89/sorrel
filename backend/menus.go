package main

import (
	"encoding/json"
	log "log/slog"
	"net/http"
	"strconv"
	"time"
)

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

	rows, err := db.Query("SELECT id, name, description, menu_section FROM items WHERE menu_id = ?", menuId)
	if err != nil {
		log.Error("error fetching menu items", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var itemId int
		var name, description string
		var section *string
		if err := rows.Scan(&itemId, &name, &description, &section); err != nil {
			log.Error("error reading row for menu item", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		menuItems = append(menuItems, MenuItem{Name: name, Description: description, Id: itemId, Section: section})
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

	rawItemId := r.PathValue("itemId")
	itemId, err := strconv.Atoi(rawItemId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var name, description, section string

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
		case "section":
			section = formValue
		default:
			http.Error(w, "extra value provided: "+formKey, http.StatusBadRequest)
			return
		}
	}

	w.Header().Add("Access-Control-Allow-Origin", "*")

	if _, err := db.Exec("UPDATE items SET name = ?, description = ?, menu_section = ? WHERE id = ?", name, description, section, itemId); err != nil {
		log.Error("error editing menu item", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleDeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	rawItemId := r.PathValue("itemId")
	itemId, err := strconv.Atoi(rawItemId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Access-Control-Allow-Origin", "*")

	_, err = db.Exec("DELETE FROM items WHERE id = ?", itemId)
	if err != nil {
		log.Error("Error deleting menu item", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleEditMenu(w http.ResponseWriter, r *http.Request) {
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

	var name string
	var date time.Time

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

	_, err = db.Exec("UPDATE menus SET name = ?, date = ? WHERE id = ?", name, date, menuId)
	if err != nil {
		log.Error("Error updating menu", "error", err)
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
	w.Header().Add("Access-Control-Allow-Origin", "*")

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
