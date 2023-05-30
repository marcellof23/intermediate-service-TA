package model

type Subscriber struct {
	Base
	ID                 int64
	GoogleSubscriberID string
	InUsed             bool
}

type TotalInUsed struct {
	NumberInused int64
	Total        int64
}
