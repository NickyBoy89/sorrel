# Sorrel: A recipe planner to share menus with friends

This project started as a fun idea: what if I could share what I'm making for dinner with friends/family in a way that's amusingly professional-looking, and work on some of my full-stack skills at the same time?

Sorrel is a self-hosted PWA application that can be easily installed on iOS or Android, and provides a web interface to edit and configure "menus". It provides both a way for clients to install the app on their own, given just an "invitation code" to avoid having to remember a username+password, and an authenticated backend for an administrator to create menus item-by-item, and share them via push notifications.

Features include:

* [x] Creating, modifying, and deleting menus
* [x] Creating, modifying, and deleting menu items
* [x] A live preview of the menu as it is being created
* [x] Ability to send a push notification when a menu has been created
* [ ] Probably a lot more...

> [!WARNING]
> Please note that this project is built for my own use, and therefore may contain implicit steps for it to work. I have tried to document as many of these as I can, but these will most likely be incomplete. If you find a missing dependency, feel free to open an issue/PR, and/or if you find something that would be generally useful for the functionality of the app, please contribute!

## Requirements

* Go 1.24+ with CGo (could probably work with earlier versions)
* Yarn
* Vite
* Sveltekit

## Building

1. Install the node requirements
    * `yarn install`
2. Build the frontend
    * `yarn run build`
3. Build the backend
    * `cd backend/ && go build`
4. Run the server binary
    * `./backend serve` (on Linux/MacOS)
    * `.\backend.exe serve` (on Windows)
5. Done! The server runs by default on port [`9031`](https://github.com/NickyBoy89/sorrel/blob/master/backend/server.go#L21)

## Architecture

The project tree is separated into the frontend and backend:
```
backend/
src/
```

* `backend/` contains the Go backend, which handles incoming http requests and persists them to an embedded SQLite database, `recipes.db`
* `src/` contains the Svelte frontend, which is compiled to a static site in `build/`, and served by the Go backend

## Configuration

By default, `sorrel` places its config file [`sorrel-config.json`](https://github.com/NickyBoy89/sorrel/blob/master/backend/main.go#L7) alongside the binary. If no config file is found, a new one is created.

> [!TIP]
> When running the CLI (`./backend`) the built-in help text (ex: `./backend --help` or `./backend serve --help`) can give information about more options available to the user

### Configuration Values in `sorrel-config.json`

* `public_key` is the VAPID public key used for push notifications
* `private_key` is the corresponding VAPID private key for push notifications
* `keycloak_hostname` (if using the default keycloak auth handler) is the url of the keycloak instance (ex: https://keycloak.my-domain.com)
* `keycloak_realm_name` (if using the default keycloak auth handler) is the keycloak realm name that the application has been configured to use
