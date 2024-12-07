// Package models is a collection of models
package models

import "github.com/volatiletech/null/v9"

// User model contains user data like GUID, email, IP and refresh token
type User struct {
	GUID    string      `json:"guid" db:"guid"`
	Email   string      `json:"email" db:"email"`
	IP      string      `json:"ip" db:"ip"`
	Refresh null.String `json:"refresh_token" db:"refresh_token"`
	Public  null.String `json:"public_key" db:"public_key"`
}
