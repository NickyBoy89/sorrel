package grocerylists

import (
	"database/sql"
	"errors"

	"github.com/NickyBoy89/sorrel/backend/internal/db"
)

func InitializeDB(db *sql.DB) error {
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS grocery_lists(
		id INTEGER PRIMARY KEY AUTOINCREMENT
	);
	`); err != nil {
		return err
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS grocery_list_contents(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		grocery_list_id INTEGER NOT NULL,
		ingredient_id INTEGER NOT NULL,
		checked BOOL NOT NULL DEFAULT FALSE,
		quantity TEXT NOT NULL,
		FOREIGN KEY(grocery_list_id) REFERENCES grocery_lists(id),
		FOREIGN KEY(ingredient_id) REFERENCES ingredients(id)
	);
	`); err != nil {
		return err
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS shared_grocery_lists(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		shopping_list_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		FOREIGN KEY(shopping_list_id) REFERENCES shopping_lists(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	);
	`); err != nil {
		return err
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS ingredients(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		category TEXT
	);
	`); err != nil {
		return err
	}

	return nil
}

func InsertGroceryList(db *sql.DB) (int64, error) {
	if resp, err := db.Exec("INSERT INTO grocery_lists DEFAULT VALUES"); err != nil {
		return 0, err
	} else {
		return resp.LastInsertId()
	}
}

func FindGroceryList(groceryListId int, db *sql.DB) ([]GroceryItem, error) {
	var count int
	if err := db.QueryRow("SELECT 1 FROM grocery_lists WHERE id = ?", groceryListId).Scan(&count); err != nil {
		return nil, err
	}

	itemRows, err := db.Query(`SELECT
			grocery_list_contents.id as id,
			ingredients.name AS name,
			ingredients.category AS category,
			quantity,
			checked
FROM
			grocery_list_contents
INNER JOIN ingredients ON
			ingredients.id = grocery_list_contents.id
WHERE
			grocery_list_contents.grocery_list_id = ?`, groceryListId)
	if err != nil {
		return nil, err
	}

	items := []GroceryItem{}

	for itemRows.Next() {
		var item GroceryItem
		if err := itemRows.Scan(&item.Id, &item.Name, &item.Category, &item.Quantity, &item.Checked); err != nil {
			return nil, errors.New("could not read item from grocery list")
		}

		items = append(items, item)
	}

	return items, nil
}

func InsertGroceryListItem(input UpdateGroceryItem, groceryListId int, db db.DatabaseLike) (int64, error) {
	if resp, err := db.Exec("INSERT INTO grocery_list_contents (grocery_list_id, ingredient_id, quantity) VALUES (?, ?, ?)", groceryListId, input.Id, input.Quantity); err != nil {
		return 0, err
	} else {
		return resp.LastInsertId()
	}
}

func UpdateGroceryList(items []UpdateGroceryItem, groceryListId int, db *sql.DB) error {
	for _, item := range items {
		if _, err := db.Exec("INSERT INTO grocery_list_contents (grocery_list_id, ingredient_id, quantity, checked) VALUES (?, ?, ?, ?)", groceryListId, item.Id, item.Quantity, item.Checked); err != nil {
			return err
		}
	}

	return nil
}

func InsertIngredient(val UpdateIngredient, db *sql.DB) (int64, error) {
	if res, err := db.Exec("INSERT INTO ingredients (name, category) VALUES (?, ?)", val.Name, val.Category); err != nil {
		return 0, err
	} else {
		return res.LastInsertId()
	}
}

func FindIngredient(ingredientId int, db *sql.DB) (Ingredient, error) {
	var ingredient Ingredient

	if err := db.QueryRow("SELECT name, category FROM ingredients WHERE id = ?", ingredientId).Scan(&ingredient.Name, &ingredient.Category); err != nil {
		return ingredient, err
	} else {
		return ingredient, nil
	}
}
