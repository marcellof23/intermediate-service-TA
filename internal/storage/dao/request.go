package dao

import "github.com/intermediate-service-ta/internal/model"

type Request struct {
	Base
	ID                   int64  `gorm:"varchar(255);primary_key;not null"`
	RequestID            string `gorm:"varchar(255);not null"`
	TotalCommand         int    `gorm:"not null"`
	TotalSuccessExecuted int    `gorm:"not null"`
	Status               model.Status
}
