-- name: GetRecipe :one
SELECT * FROM recipes WHERE id = ?;

-- name: GetRecipeIngredients :many
SELECT * FROM recipe_ingredients WHERE recipe_id = ?;

-- name: FindRecipes :many
SELECT * FROM recipes ORDER BY id;
