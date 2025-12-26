package grocerylists

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	log "log/slog"
	"net/http"
	"strconv"

	"github.com/NickyBoy89/sorrel/backend/internal/db"
)

type GroceryItem struct {
	// The id of the item in the grocery list
	Id int `json:"id"`
	// The name of the item
	Name string `json:"name"`
	// The general category in the grocery store where this item is located
	Category *string `json:"category,omitempty"`
	// A description of how much of this item is required
	Quantity string `json:"quantity"`
	// Whether this item has been marked as complete
	Checked bool `json:"checked"`
}

type CreateGroceryItem struct {
	Quantity string `json:"quantity"`
	Checked  bool   `json:"checked"`
}

type UpdateGroceryItem struct {
	Id       int    `json:"id"`
	Quantity string `json:"quantity"`
	Checked  bool   `json:"checked"`
}

// Represents a single Ingredient in the library of ingredients
type Ingredient struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Category *string `json:"category,omitempty"`
}

type UpdateIngredient struct {
	Name     string  `json:"name"`
	Category *string `json:"category,omitempty"`
}

func RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/grocery_list/{id}", handleGroceryListAction)
	mux.HandleFunc("/api/v1/grocery_list", handleCreateGroceryList)

	mux.HandleFunc("/api/v1/ingredients/{id}", handleIngredientAction)
	mux.HandleFunc("/api/v1/ingredients", handleCreateIngredient)
}

// `handleGroceryListAction` handles the CRUD actions of a grocery list
func handleGroceryListAction(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	groceryId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch r.Method {
	// Get a grocery list
	case http.MethodGet:
		// Find what items are on the list

		items, err := FindGroceryList(groceryId, db.DB)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			log.Error("error querying grocery items", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		if err := json.NewEncoder(w).Encode(items); err != nil {
			log.Error("error encoding grocery items", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	// Edit a shopping list
	case http.MethodPatch:
		defer r.Body.Close()

		var groceryItems []UpdateGroceryItem
		if err := json.NewDecoder(r.Body).Decode(&groceryItems); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// tx, err := db.DB.Begin()
		// if err != nil {
		// 	log.Error("error beginning transaction", "error", err)
		// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		// 	return
		// }
		// defer tx.Rollback()

		if err := UpdateGroceryList(groceryItems, groceryId, db.DB); err != nil {
			log.Error("Error updating grocery list", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// if err := tx.Commit(); err != nil {
		// 	log.Error("error committing grocery list change", "error", err)
		// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		// 	return
		// }

	// Delete a shopping list
	case http.MethodDelete:
		if _, err := db.DB.Exec("DELETE FROM grocery_lists WHERE id = ?", groceryId); err != nil {
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
	w.Header().Add("Access-Control-Allow-Origin", "*")

	if groceryListId, err := InsertGroceryList(db.DB); err != nil {
		log.Error("error while creating new grocery list", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else {
		fmt.Fprintf(w, "%d", groceryListId)
	}
}

func handleIngredientAction(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	ingredientId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		ing, err := FindIngredient(ingredientId, db.DB)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			log.Error("error while finding ingredient", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(ing); err != nil {
			log.Error("error encoding ingredient", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func handleCreateIngredient(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	defer r.Body.Close()

	var item UpdateIngredient

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if id, err := InsertIngredient(item, db.DB); err != nil {
		log.Error("error inserting ingredient", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else {
		fmt.Fprintf(w, "%d", id)
	}
}
