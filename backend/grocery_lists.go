package main

import (
	"encoding/json"
	log "log/slog"
	"net/http"
	"strconv"
)

type GroceryItem struct {
	Name     string  `json:"name"`
	Category *string `json:"category,omitempty"`
	Quantity string  `json:"quantity"`
	Checked  bool    `json:"checked"`
}

// `handleGroceryListAction` handles the CRUD actions of a grocery list
func handleGroceryListAction(w http.ResponseWriter, r *http.Request) {
	rawGroceryId := r.PathValue("groceryListId")

	groceryId, err := strconv.Atoi(rawGroceryId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch r.Method {
	// Get a grocery list
	case http.MethodGet:
		// Find what items are on the list
		groceryListItems, err := db.Query(`SELECT
	name,
	quantity,
	category
FROM
	grocery_items
INNER JOIN grocery_list_contents ON 
	grocery_list_contents.grocery_item_id = grocery_items.id
WHERE grocery_list_contents.grocery_list_id = ?`, groceryId)
		if err != nil {
			log.Error("error querying grocery items", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		groceryItems := []GroceryItem{}

		for groceryListItems.Next() {
			// Read the grocery item
			var item GroceryItem
			if err := groceryListItems.Scan(
				&item.Name,
				&item.Quantity,
				&item.Category,
			); err != nil {
				log.Error("error reading grocery item", "error", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			groceryItems = append(groceryItems, item)
		}

		if err := json.NewEncoder(w).Encode(groceryItems); err != nil {
			log.Error("error encoding grocery items", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	// Edit a shopping list
	case http.MethodPatch:
		defer r.Body.Close()
		log.Error("Unimplemented")
		http.Error(w, "Unimplemented", http.StatusNotFound)
		return
	// Delete a shopping list
	case http.MethodDelete:
		if _, err := db.Exec("DELETE FROM grocery_lists WHERE id = ?", groceryId); err != nil {
			log.Error("error deleting grocery list", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func handleCreateGroceryList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		log.Error("unimplemented")
		http.Error(w, "unimplemented", http.StatusNotFound)
		return
	case http.MethodPost:
		if _, err := db.Exec("INSERT INTO grocery_lists DEFAULT VALUES"); err != nil {
			log.Error("error creating new grocery item", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}
