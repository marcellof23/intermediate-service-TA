package model

import "github.com/google/uuid"

type File struct {
	Base
	ID           uuid.UUID
	Filename     string
	OriginalName string
	Mimetype     string
	Url          string
	Client       string
	Size         int64
}
