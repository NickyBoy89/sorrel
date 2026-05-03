CREATE TABLE recipes (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL DEFAULT 'Unnamed Recipe',
  description TEXT,
  color TEXT,
  yield TEXT,
  active_time TEXT,
  total_time TEXT
);

CREATE TABLE recipe_ingredients (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  recipe_id INTEGER NOT NULL,
  ingredient_id INTEGER NOT NULL,
  FOREIGN KEY(recipe_id) REFERENCES recipes(id),
  FOREIGN KEY(ingredient_id) REFERENCES ingredients(id)
);

CREATE TABLE recipe_step (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  recipe_id INTEGER NOT NULL,
  sequence_number INTEGER NOT NULL,
  FOREIGN KEY(recipe_id) REFERENCES recipes(id)
);

-- name: GetRecipe :one
SELECT * FROM recipes WHERE id = ?;

-- name: GetRecipeIngredients :many
SELECT * FROM recipe_ingredients WHERE recipe_id = ?;

-- name: FindRecipes :many
SELECT * FROM recipes ORDER BY id;
