package grocerylists

import (
	"encoding/json"
	"testing"

	"github.com/NickyBoy89/sorrel/backend/internal/db"
	"github.com/NickyBoy89/sorrel/backend/internal/testutils"
	_ "github.com/mattn/go-sqlite3"
)

func TestCreateGroceryListV1(t *testing.T) {
	f, err := testutils.NewMockFramework(RegisterHandlers, InitializeDB)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	f.CheckedPost(t, "/api/v1/grocery_list", nil)

	var count int
	if err := db.DB.QueryRow("SELECT COUNT(*) FROM grocery_lists").Scan(&count); err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Error("expected exactly one grocery list", "length", count)
	}
}

func addGroceryItem(t *testing.T, name string) int64 {
	res, err := InsertIngredient(UpdateIngredient{Name: name}, db.DB)

	if err != nil {
		t.Error(err)
	}

	return res
}

func TestGetGroceryListItemsV1(t *testing.T) {
	f, err := testutils.NewMockFramework(RegisterHandlers, InitializeDB)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	f.CheckedPost(t, "/api/v1/grocery_list", nil)

	for _, fruitId := range []int64{
		addGroceryItem(t, "Apple"),
		addGroceryItem(t, "Banana"),
		addGroceryItem(t, "Grapefruit"),
	} {
		if _, err := InsertGroceryListItem(CreateGroceryItem{Quantity: "1 fruit", Checked: false}, int(fruitId), 1, db.DB); err != nil {
			t.Error(err)
		}
	}

	addGroceryItem(t, "Pear")

	resp := f.CheckedGet(t, "/api/v1/grocery_list/1")
	defer resp.Body.Close()

	var groceryItems []GroceryItem
	if err := json.NewDecoder(resp.Body).Decode(&groceryItems); err != nil {
		t.Fatal(err)
	}

	if len(groceryItems) != 3 {
		t.Error("Expected 3 grocery items, got", len(groceryItems), "structure", groceryItems)
	}
}
