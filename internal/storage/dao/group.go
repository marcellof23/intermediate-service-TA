package dao

type Group struct {
	Base
	ID        int64  // Group ID
	GroupName string `gorm:"varchar(255);not null"`
}
