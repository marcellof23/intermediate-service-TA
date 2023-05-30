package model

type Role string

const (
	Admin  Role = "Admin"
	Normal Role = "Normal"
)

// User represents a Unix user
type User struct {
	Base
	ID           int64  // User ID
	Username     string // Username of the user
	Password     string `json:"-"` // Password of the user
	SubscriberID string
	Role         Role  // Role of the user
	GroupID      int64 // Group ID
}

// Create a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
