package grocerylists

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
		if _, err := InsertGroceryListItem(UpdateGroceryItem{Quantity: "1 fruit", Checked: false, Id: int(fruitId)}, 1, db.DB); err != nil {
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

func TestDeleteGroceryList(t *testing.T) {
	f, err := testutils.NewMockFramework(RegisterHandlers, InitializeDB)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	resp := f.CheckedPost(t, "/api/v1/grocery_list", nil)

	var listId int
	if _, err := fmt.Fscan(resp.Body, &listId); err != nil {
		t.Fatal(err)
	}

	f.CheckedGet(t, fmt.Sprintf("/api/v1/grocery_list/%d", listId))

	f.CheckedRequest(t, http.MethodDelete, fmt.Sprintf("/api/v1/grocery_list/%d", listId), nil)

	resp = f.Get(t, fmt.Sprintf("/api/v1/grocery_list/%d", listId))

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("actual code did not equal expected, actual %d, expected %d", resp.StatusCode, http.StatusNotFound)
	}
}

func createIngredient(t *testing.T, f *testutils.MockFramework, ingredient UpdateIngredient) int {

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(ingredient); err != nil {
		t.Fatal(err)
	}

	return testutils.IdOf(t, f.CheckedPost(t, "/api/v1/ingredients", &buf))
}

func TestAddItemsToGroceryList(t *testing.T) {
	f, err := testutils.NewMockFramework(RegisterHandlers, InitializeDB)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	listId := testutils.IdOf(t, f.CheckedPost(t, "/api/v1/grocery_list", nil))

	eggId := createIngredient(t, f, UpdateIngredient{Name: "Egg"})
	milkId := createIngredient(t, f, UpdateIngredient{Name: "Milk"})

	if initial, err := json.Marshal([]UpdateGroceryItem{{Id: eggId}, {Id: milkId}}); err != nil {
		t.Fatal(err)
	} else {
		f.CheckedRequest(t, http.MethodPatch, fmt.Sprintf("/api/v1/grocery_list/%d", listId), bytes.NewBuffer(initial))
	}

	resp := f.CheckedGet(t, fmt.Sprintf("/api/v1/grocery_list/%d", listId))
	defer resp.Body.Close()

	insertedItems := []GroceryItem{}
	if err := json.NewDecoder(resp.Body).Decode(&insertedItems); err != nil {
		t.Fatal(err)
	}

	if len(insertedItems) != 2 {
		t.Errorf("expected two items, got %d", len(insertedItems))
	}
}

func TestAddAnotherItemToList(t *testing.T) {
	f, err := testutils.NewMockFramework(RegisterHandlers, InitializeDB)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	listId := testutils.IdOf(t, f.CheckedPost(t, "/api/v1/grocery_list", nil))

	eggId := createIngredient(t, f, UpdateIngredient{Name: "Egg"})
	milkId := createIngredient(t, f, UpdateIngredient{Name: "Milk"})
	cornId := createIngredient(t, f, UpdateIngredient{Name: "Corn stalks"})

	if encoded, err := json.Marshal([]UpdateGroceryItem{{Id: eggId}, {Id: milkId}}); err != nil {
		t.Fatal(err)
	} else {
		f.CheckedRequest(t, http.MethodPatch, fmt.Sprintf("/api/v1/grocery_list/%d", listId), bytes.NewBuffer(encoded))
	}

	if encoded, err := json.Marshal([]UpdateGroceryItem{{Id: eggId}, {Id: milkId}, {Id: cornId}}); err != nil {
		t.Fatal(err)
	} else {
		f.CheckedRequest(t, http.MethodPatch, fmt.Sprintf("/api/v1/grocery_list/%d", listId), bytes.NewBuffer(encoded))
	}

	resp := f.CheckedGet(t, fmt.Sprintf("/api/v1/grocery_list/%d", listId))
	defer resp.Body.Close()

	updatedItems := []GroceryItem{}
	if err := json.NewDecoder(resp.Body).Decode(&updatedItems); err != nil {
		t.Fatal(err)
	}

	if len(updatedItems) != 3 {
		t.Errorf("expected three items, got %d", len(updatedItems))
	}
}
