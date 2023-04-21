package dao

import (
	"github.com/intermediate-service-ta/internal/model"
)

// User represents a Unix user
type User struct {
	Base
	ID       int64      // User ID
	Username string     `gorm:"varchar(255);not null"`
	Password string     `gorm:"varchar(255);not null"`
	Role     model.Role `gorm:"varchar(255);not null"`
	GroupID  int64      // Group ID
}

func ToUserDTO(u User) model.User {
	return model.User{
		ID:       u.ID,
		Username: u.Username,
		Password: u.Password,
		Role:     u.Role,
		GroupID:  u.GroupID,
		Base: model.Base{
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
	}
}

func ToUserDAO(u model.User) User {
	return User{
		ID:       u.ID,
		Username: u.Username,
		Password: u.Password,
		Role:     u.Role,
		GroupID:  u.GroupID,
		Base: Base{
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
	}
}
