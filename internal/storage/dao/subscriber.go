package dao

type Subscriber struct {
	Base
	ID                 int64  `gorm:"varchar(255);primary_key;not null"`
	GoogleSubscriberID string `gorm:"varchar(255);not null"`
	InUsed             bool
}
