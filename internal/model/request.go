package model

const (
	Pending Status = "Pending"
	Failed  Status = "Failed"
	Success Status = "Success"
)

type Status string

type Request struct {
	Base
	ID                   int64
	RequestID            string
	TotalCommand         int
	TotalSuccessExecuted int
	Status               Status
}

type RequestCommand struct {
	RequestID    string
	TotalCommand int
}
