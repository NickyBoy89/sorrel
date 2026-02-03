package recipes

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	grocerylists "github.com/NickyBoy89/sorrel/backend/grocery_lists"
	"github.com/NickyBoy89/sorrel/backend/internal/db"
	"github.com/NickyBoy89/sorrel/backend/internal/middleware"
)

type Recipe struct {
	Id          RecipeId                   `json:"id"`
	Name        string                     `json:"name"`
	Color       string                     `json:"color"`
	Ingredients []grocerylists.GroceryItem `json:"ingredients"`
}

func RegisterHandlers(mux *http.ServeMux, auth middleware.AuthMiddleware) {
	mux.HandleFunc("/api/v1/recipes", handleListRecipes)
	mux.Handle("/api/v1/recipes/{id}", auth.RequireAuth(http.HandlerFunc(handleRecipeAction)))
}

func handleListRecipes(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "*")

	switch r.Method {
	case http.MethodGet:
		recipeIds, err := FindRecipeIds(db.DB)
		if err != nil {
			slog.Error("error reading recipe ids", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(recipeIds); err != nil {
			slog.Error("error encoding recipe ids", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		tx, err := db.DB.Begin()
		if err != nil {
			slog.Error("error starting reansaction to fetch recipe", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		defer r.Body.Close()

		var input CreateRecipeInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		slog.Debug("creating new recipe", "input", input)

		recipe, err := InsertRecipe(input, tx)
		if err != nil {
			slog.Error("error creating recipe", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(recipe); err != nil {
			slog.Error("error encoding recipe", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(); err != nil {
			slog.Error("error commiting recipe read", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

func handleRecipeAction(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "*")

	recipeId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		slog.Error("error starting reansaction to fetch recipe", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	switch r.Method {
	case http.MethodGet:
		recipe, err := FindRecipe(RecipeId(recipeId), tx)
		if err != nil {
			slog.Error("error fetching recipe", "error", err, "recipeId", recipeId)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(recipe); err != nil {
			slog.Error("error encoding recipe", "error", err, "recipeId", recipeId)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodPatch:
		var input map[string]any
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if updatedName, ok := input["name"]; ok {
			if _, err := db.DB.Exec("UPDATE recipes SET name = ? WHERE id = ?", updatedName, recipeId); err != nil {
				slog.Error("error updating name of recipe", "error", err, "recipeId", recipeId)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		slog.Error("error commiting recipe read", "error", err, "recipeId", recipeId)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
