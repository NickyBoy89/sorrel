package auth

import "time"

// UserInfo contains the essential user information we want to cache
type UserInfo struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	// Add other user fields you need
}

// SessionData represents the data we'll store for each session
type SessionData struct {
	AccessToken string    `json:"access_token"`
	UserInfo    UserInfo  `json:"user_info"`
	CreatedAt   time.Time `json:"created_at"`
}
