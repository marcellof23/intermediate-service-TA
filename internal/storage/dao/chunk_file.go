package dao

import "github.com/intermediate-service-ta/internal/model"

type ChunkFile struct {
	Base
	ID       int64  `gorm:"varchar(255);primary_key;not null"`
	FileID   int64  `gorm:"varchar(255);not null"`
	Filename string `gorm:"varchar(255);not null"`
	Order    int    `gorm:"not null"`
	Size     int64  `gorm:"not null"`
	Client   string `gorm:"varchar(255);not null"`
}

func ToFileChunkDAO(fc *model.ChunkFile) *ChunkFile {
	return &ChunkFile{
		ID:       fc.ID,
		Filename: fc.Filename,
		FileID:   fc.FileID,
		Order:    fc.Order,
		Size:     fc.Size,
		Client:   fc.Client,
		Base: Base{
			CreatedAt: fc.CreatedAt,
			UpdatedAt: fc.UpdatedAt,
		},
	}
}

func ToFileChunkDTO(fc ChunkFile) model.ChunkFile {
	return model.ChunkFile{
		ID:       fc.ID,
		Filename: fc.Filename,
		FileID:   fc.FileID,
		Order:    fc.Order,
		Size:     fc.Size,
		Client:   fc.Client,
		Base: model.Base{
			CreatedAt: fc.CreatedAt,
			UpdatedAt: fc.UpdatedAt,
		},
	}
}
