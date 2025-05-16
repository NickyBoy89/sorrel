package main

import (
	"encoding/json"
	log "log/slog"
	"net/http"
	"strconv"
)

type GroceryItem struct {
	Id       int     `json:"id"`
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
	grocery_items.id AS id,
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
				&item.Id,
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

		var groceryItemIds []int
		if err := json.NewDecoder(r.Body).Decode(&groceryItemIds); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Error("error beginning transaction", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		// Delete all current mappings
		if _, err := tx.Exec("DELETE FROM grocery_list_contents WHERE grocery_list_id = ?", groceryId); err != nil {
			log.Error("error deleting items from grocery list", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		insertItemQuery, err := tx.Prepare("INSERT INTO grocery_list_contents (grocery_list_id, grocery_item_id) VALUES (?, ?)")
		if err != nil {
			log.Error("error preparing query to add grocery item", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer insertItemQuery.Close()

		// Create new ones from the ids provided
		for _, itemId := range groceryItemIds {
			if _, err := insertItemQuery.Exec(groceryId, itemId); err != nil {
				log.Error("error adding new items to grocery list", "error", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		if err := tx.Commit(); err != nil {
			log.Error("error committing grocery list change", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

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
