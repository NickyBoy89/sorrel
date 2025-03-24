package main

import "database/sql"

func initDb(db *sql.DB) error {
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS menus(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		date TIMESTAMP
	);
	`); err != nil {
		return err
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS items(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		menu_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		description TEXT,
		menu_section TEXT,
		FOREIGN KEY(menu_id) REFERENCES menus(id)
	);
	`); err != nil {
		return err
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		display_name TEXT NOT NULL
	);
	`); err != nil {
		return err
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS notification_subscriptions(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		endpoint TEXT NOT NULL,
		keys_auth TEXT NOT NULL,
		keys_p256dh TEXT NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);
	`); err != nil {
		return err
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS shared_menus(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		menu_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		FOREIGN KEY(menu_id) REFERENCES menus(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	);
	`); err != nil {
		return err
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS shopping_lists(
		id INTEGER PRIMARY KEY AUTOINCREMENT
	);
	`); err != nil {
		return err
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS shared_lists(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		shopping_list_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		FOREIGN KEY(shopping_list_id) REFERENCES shopping_lists(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	);
	`); err != nil {
		return err
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS grocery_item(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		category TEXT,
		quantity TEXT NOT NULL,
		checked BOOL NOT NULL
	);
	`); err != nil {
		return err
	}

	return nil
}
