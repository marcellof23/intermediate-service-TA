package model

type File struct {
	Base
	ID           int64
	Filename     string
	OriginalName string
	Mimetype     string
	Url          string
	Client       string
	Size         int64
}
