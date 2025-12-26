package grocerylists

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/NickyBoy89/sorrel/backend/internal/testutils"
)

func TestAddIngredient(t *testing.T) {
	f, err := testutils.NewMockFramework(RegisterHandlers, InitializeDB)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	f.CheckedPost(t, "/api/v1/ingredients", bytes.NewBufferString("{\"name\": \"Cheese\"}"))
}

func TestGetIngredient(t *testing.T) {
	f, err := testutils.NewMockFramework(RegisterHandlers, InitializeDB)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	cheese := UpdateIngredient{
		Name: "Cheese",
	}

	encoded, err := json.Marshal(cheese)
	if err != nil {
		t.Fatal(err)
	}

	resp := f.CheckedPost(t, "/api/v1/ingredients", bytes.NewBufferString(string(encoded)))

	var ingredientId int

	if _, err := fmt.Fscan(resp.Body, &ingredientId); err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	resp = f.CheckedGet(t, fmt.Sprintf("/api/v1/ingredients/%d", ingredientId))
	var fetched UpdateIngredient
	if err := json.NewDecoder(resp.Body).Decode(&fetched); err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if fetched != cheese {
		t.Errorf("actual did not equal expected, expected: %v, actual: %v", cheese, fetched)
	}

	resp = f.Get(t, "/api/v1/ingredients/1234")

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("actual code did not equal expected, actual %d, expected %d", resp.StatusCode, http.StatusNotFound)
	}
}
