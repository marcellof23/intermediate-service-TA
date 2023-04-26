package dao

import (
	"github.com/intermediate-service-ta/internal/model"
)

var MemorySlice map[string]int

type File struct {
	Base
	ID           int64  `gorm:"varchar(255);primary_key;not null"`
	Filename     string `gorm:"varchar(255);not null"`
	OriginalName string `gorm:"varchar(255);not null"`
	Size         int64  `gorm:"not null"`
	Client       string `gorm:"varchar(255);not null"`
}

func ToFileDAO(f *model.File) *File {
	return &File{
		ID:           f.ID,
		Filename:     f.Filename,
		OriginalName: f.OriginalName,
		Size:         f.Size,
		Client:       f.Client,
		Base: Base{
			CreatedAt: f.CreatedAt,
			UpdatedAt: f.UpdatedAt,
		},
	}
}

func ToFileDTO(f File) model.File {
	return model.File{
		ID:           f.ID,
		Filename:     f.Filename,
		OriginalName: f.OriginalName,
		Size:         f.Size,
		Client:       f.Client,
		Base: model.Base{
			CreatedAt: f.CreatedAt,
			UpdatedAt: f.UpdatedAt,
		},
	}
}
