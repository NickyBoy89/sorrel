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
	"github.com/NickyBoy89/sorrel/backend/internal/middleware"
	"github.com/NickyBoy89/sorrel/backend/internal/push"
)

type GroceryList struct {
	Id   int     `json:"id"`
	Name *string `json:"name,omitempty"`
}

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
	IngredientId int    `json:"id"`
	Quantity     string `json:"quantity"`
	Checked      bool   `json:"checked"`
}

type IngredientId int64

// Represents a single Ingredient in the library of ingredients
type Ingredient struct {
	Id       IngredientId `json:"id"`
	Name     string       `json:"name"`
	Category *string      `json:"category,omitempty"`
}

type UpdateIngredient struct {
	Name     string  `json:"name"`
	Category *string `json:"category,omitempty"`
}

func RegisterHandlers(mux *http.ServeMux, auth middleware.AuthMiddleware) {
	mux.Handle("/api/v1/grocery_list/{id}/item/{itemId}", auth.RequireAuth(http.HandlerFunc(handleDeleteItem)))
	mux.Handle("/api/v1/grocery_list/{id}", http.HandlerFunc(handleGroceryListAction))
	mux.Handle("/api/v1/grocery_list/{id}/share", auth.RequireAuth(http.HandlerFunc(handleShareGroceryList)))
	mux.Handle("/api/v1/grocery_list", http.HandlerFunc(handleCreateGroceryList))

	mux.Handle("/api/v1/ingredients/{id}", auth.RequireAuth(http.HandlerFunc(handleIngredientAction)))
	mux.Handle("/api/v1/ingredients", auth.RequireAuth(http.HandlerFunc(handleCreateIngredient)))
}

// `handleGroceryListAction` handles the CRUD actions of a grocery list
func handleGroceryListAction(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "*")

	groceryId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodOptions:
	case http.MethodPost:
		defer r.Body.Close()
		var item UpdateGroceryItem

		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if item.IngredientId == 0 || item.Quantity == "" {
			http.Error(w, "one of the \"id\" or \"quantity\" fields were missing in the request input", http.StatusBadRequest)
			return
		}

		if _, err := InsertGroceryListItem(item, groceryId, db.DB); err != nil {
			log.Error("error inserting new grocery item", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	case http.MethodGet: // Find what items are on the list
		items, err := FindGroceryListItems(groceryId, db.DB)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			log.Error("error querying grocery items", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
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

		if err := UpdateGroceryListItem(groceryItems, groceryId, db.DB); err != nil {
			log.Error("Error updating grocery list", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// if err := tx.Commit(); err != nil {
		// 	log.Error("error committing grocery list change", "error", err)
		// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		// 	return
		// }
	case http.MethodPut:
		defer r.Body.Close()

		var list GroceryList
		if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := UpdateGroceryList(list, db.DB); err != nil {
			log.Error("error while updating grocery list", "error", err, "id", list.Id)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

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

func handleShareGroceryList(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	groceryId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		log.Error("error beginning transaction", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	defer r.Body.Close()

	users := []int{}
	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	list, err := FindGroceryList(groceryId, tx)
	if err != nil {
		log.Error("error fetching grocery list", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	listName := "Groceries"
	if list.Name != nil {
		listName = *list.Name
	}

	msg := push.PushMessage{
		Message:   fmt.Sprintf("A grocery list has been shared with you: %s", listName),
		ActionUrl: fmt.Sprintf("/grocery_lists/list/?id=%d", 1),
	}

	resp := make(map[int]bool)

	for _, userId := range users {
		if sent, err := push.SendNotificationToUser(tx, userId, msg); err != nil {
			log.Error("error sharing list with user", "error", err, "userId", userId, "groceryListId", groceryId)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			resp[userId] = sent
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error("error committing notification changes", "error", err)
		http.Error(w, "error committing notification changes", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("error encoding menu share response", "error", err)
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func handleCreateGroceryList(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "*")

	switch r.Method {
	case http.MethodGet:
		ids, err := FindGroceryLists(db.DB)
		if err != nil {
			log.Error("error while querying grocery lists", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(ids); err != nil {
			log.Error("error while encoding grocery list ids", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		if groceryListId, err := InsertGroceryList(db.DB); err != nil {
			log.Error("error while creating new grocery list", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		} else {
			fmt.Fprintf(w, "%d", groceryListId)
		}
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func handleDeleteItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "*")

	groceryListId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	groceryItemId, err := strconv.Atoi(r.PathValue("itemId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodOptions:
		return
	case http.MethodDelete:
		if err := DeleteGroceryListItem(groceryListId, groceryItemId, db.DB); err != nil {
			log.Error("error while deleting item in grocery list", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
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
		ing, err := FindIngredient(IngredientId(ingredientId), db.DB)
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

	switch r.Method {
	case http.MethodGet:
		ingredients, err := FindIngredients(db.DB)
		if err != nil {
			log.Error("error fetching ingredients", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(ingredients); err != nil {
			log.Error("error encoding ingredients", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
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
}
