package recipes

import (
	"database/sql"

	grocerylists "github.com/NickyBoy89/sorrel/backend/grocery_lists"
	"github.com/NickyBoy89/sorrel/backend/internal/db"
)

func InitializedDb(db *sql.DB) error {
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS recipe_ingredients(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	recipe_id INTEGER NOT NULL,
	ingredient_id INTEGER NOT NULL,
	FOREIGN KEY(recipe_id) REFERENCES recipes(id),
	FOREIGN KEY(ingredient_id) REFERENCES ingredients(id)
	)`); err != nil {
		return err
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS recipes(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		color TEXT
		);`); err != nil {
		return err
	}

	return nil
}

type RecipeId int64
type RecipeIngredientId int64

type CreateRecipeInput struct {
	Name          string                      `json:"name"`
	Color         string                      `json:"color"`
	IngredientIds []grocerylists.IngredientId `json:"ingredient_ids"`
}

func FindRecipeIds(db *sql.DB) ([]RecipeId, error) {
	recipeIds := []RecipeId{}

	recipeIdRows, err := db.Query("SELECT id FROM recipes")
	if err != nil {
		return nil, err
	}

	for recipeIdRows.Next() {
		var recipeId RecipeId
		if err := recipeIdRows.Scan(&recipeId); err != nil {
			return nil, err
		}

		recipeIds = append(recipeIds, recipeId)
	}

	return recipeIds, nil
}

func FindRecipe(recipeId RecipeId, db *sql.Tx) (Recipe, error) {
	recipe := Recipe{
		Id:          recipeId,
		Ingredients: []grocerylists.GroceryItem{},
	}

	if err := db.QueryRow("SELECT name, color FROM recipes WHERE id = ?", recipeId).Scan(&recipe.Name, &recipe.Color); err != nil {
		return recipe, err
	}

	ingredientRows, err := db.Query(`SELECT
		recipe_ingredients.id AS id,
		ingredients.name AS name,
		ingredients.category AS category
	FROM
		recipe_ingredients
	INNER JOIN ingredients ON
		ingredients.id = recipe_ingredients.ingredient_id
	WHERE
		recipe_ingredients.recipe_id = ?`, recipeId)
	if err != nil {
		return recipe, err
	}

	for ingredientRows.Next() {
		var ingredient grocerylists.GroceryItem
		if err := ingredientRows.Scan(&ingredient.Id, &ingredient.Name, &ingredient.Category); err != nil {
			return recipe, err
		}

		recipe.Ingredients = append(recipe.Ingredients, ingredient)
	}

	return recipe, nil
}

func InsertRecipe(input CreateRecipeInput, db *sql.Tx) (RecipeId, error) {
	res, err := db.Exec("INSERT INTO recipes (name, color) VALUES (?, ?)", input.Name, input.Color)
	if err != nil {
		return 0, err
	}

	recipeId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	for _, ingredientId := range input.IngredientIds {
		if err := InsertRecipeIngredient(ingredientId, RecipeId(recipeId), db); err != nil {
			return 0, err
		}
	}

	return RecipeId(recipeId), nil
}

func InsertRecipeIngredient(ingredientId grocerylists.IngredientId, recipeId RecipeId, db db.DatabaseLike) error {
	if _, err := db.Exec("INSERT INTO recipe_ingredients (recipe_id, ingredient_id) VALUES (?, ?)", recipeId, ingredientId); err != nil {
		return err
	}

	return nil
}
