package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	log "log/slog"
	"net/http"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

const dbTestFile = "test.db"

func setupGroceryTestingDB() {
	if err := os.Remove(dbTestFile); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			panic(err)
		}
	}

	testDB, err := sql.Open("sqlite3", dbTestFile)
	if err != nil {
		panic(err)
	}

	initDb(testDB)

	db = testDB
}

func initHandlers() {
	SetAuthMethod("none")

	// Application
	http.HandleFunc("/api/menu/{menuId}", handleGetMenu)
	http.Handle("/api/menu/{menuId}/edit", authHandlerFunc(handleEditMenu))
	http.HandleFunc("/api/menu/{menuId}/items", handleGetMenuItems)
	http.Handle("/api/menu/{menuId}/create-item", authHandlerFunc(handleCreateMenuItem))
	http.HandleFunc("/api/menu/list", handleListMenus)
	http.Handle("/api/menu/create", authHandlerFunc(handleCreateMenu))
	http.Handle("/api/menu/share", authHandlerFunc(handleShareMenu))
	http.Handle("/api/items/{itemId}/edit", authHandlerFunc(handleEditMenuItem))
	http.Handle("/api/items/{itemId}/delete", authHandlerFunc(handleDeleteMenuItem))

	// Grocery lists
	http.HandleFunc("/api/v1/grocery_list/{groceryListId}", handleGroceryListAction)
	http.HandleFunc("/api/v1/grocery_list", handleCreateGroceryList)

	// User
	http.HandleFunc("/api/validate-id", handleCheckUserId)
	http.Handle("/api/users", authHandlerFunc(handleListUsers))
	http.Handle("/api/user", authHandlerFunc(handleGetUser))

	// Push
	http.HandleFunc("/api/push/public-key", handleVAPIDPublicKeyRequest)
	http.HandleFunc("/api/push/subscribe", handlePushSubscription)
	http.HandleFunc("/api/push/validate", handleCheckSubscription)
}

// Asynchronusly starts a server and returns a handle to it
func startHTTPServer() *http.Server {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", serverPort),
	}

	if authHandlerFunc == nil {
		initHandlers()
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Error("error shutting down http server", "error", err)
			return
		}
	}()

	return srv
}

func localURL(relativeUrl string) string {
	url := fmt.Sprintf("http://localhost:%d%s", serverPort, relativeUrl)
	return url
}

func localAPIRequest(t *testing.T, method string, relativeUrl string) *http.Response {
	req, err := http.NewRequest(method, localURL(relativeUrl), nil)
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	return resp
}

func TestCreateGroceryListV1(t *testing.T) {
	setupGroceryTestingDB()
	s := startHTTPServer()
	defer s.Shutdown(context.Background())

	resp := localAPIRequest(t, http.MethodPost, "/api/v1/grocery_list")

	if resp.StatusCode != 200 {
		data, _ := io.ReadAll(resp.Body)
		t.Error("Response was not successful", "code", resp.StatusCode, "body", string(data))
	}

	rows, err := db.Query("SELECT * FROM grocery_lists")
	if err != nil {
		t.Fatal(err)
	}

	groceryListIds := []int{}

	for rows.Next() {
		var listId int
		if err := rows.Scan(&listId); err != nil {
			t.Fatal(err)
		}

		groceryListIds = append(groceryListIds, listId)
	}

	if len(groceryListIds) != 1 {
		t.Error("expected 1 grocery list, got", len(groceryListIds), "ids:", groceryListIds)
	}
}

func addGroceryItem(t *testing.T, name, quantity string) int {
	if resp, err := db.Exec("INSERT INTO grocery_items (name, quantity) VALUES (?, ?)", name, quantity); err != nil {
		t.Fatal(err)
	} else {
		if id, err := resp.LastInsertId(); err != nil {
			t.Fatal(err)
		} else {
			return int(id)
		}
	}

	panic("Unreachable")
}

func addMapping(t *testing.T, groceryListId, itemId int) {
	if _, err := db.Exec("INSERT INTO grocery_list_contents (grocery_list_id, grocery_item_id) VALUES (?, ?)", groceryListId, itemId); err != nil {
		t.Fatal(err)
	}
}

func TestGetGroceryListItemsV1(t *testing.T) {
	setupGroceryTestingDB()
	s := startHTTPServer()
	defer s.Shutdown(context.Background())

	resp := localAPIRequest(t, http.MethodPost, "/api/v1/grocery_list")

	if resp.StatusCode != http.StatusOK {
		t.Fatal(resp.StatusCode)
	}

	a := addGroceryItem(t, "Apple", "1")
	b := addGroceryItem(t, "Banana", "2")
	c := addGroceryItem(t, "Grapefruit", "3")
	addGroceryItem(t, "Pear", "4")

	addMapping(t, 1, a)
	addMapping(t, 1, b)
	addMapping(t, 1, c)

	resp = localAPIRequest(t, http.MethodGet, "/api/v1/grocery_list/1")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatal("error from server", "code", resp.StatusCode)
	}

	var groceryItems []GroceryItem
	if err := json.NewDecoder(resp.Body).Decode(&groceryItems); err != nil {
		t.Fatal(err)
	}

	if len(groceryItems) != 3 {
		t.Error("Expected 3 grocery items, got", len(groceryItems), "structure", groceryItems)
	}
}
