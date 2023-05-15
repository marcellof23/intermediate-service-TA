package model

type ChunkFile struct {
	Base
	ID       int64
	FileID   int64
	Filename string
	Order    int
	Size     int64
	Client   string
}
