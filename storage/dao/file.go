package dao

import (
	"github.com/google/uuid"

	"github.com/intermediate-service-ta/model"
)

type File struct {
	Base
	ID           uuid.UUID `gorm:"varchar(255);primary_key;not null"`
	Filename     string    `gorm:"varchar(255);not null"`
	OriginalName string    `gorm:"varchar(255);not null"`
	Mimetype     string    `gorm:"varchar(255);not null"`
	Url          string    `gorm:"varchar(255);not null"`
	Size         int64     `gorm:"type:int;not null"`
	Extension    string    `gorm:"varchar(255);not null"`
}

func ToFileDAO(f *model.File) *File {
	return &File{
		ID:           f.ID,
		Filename:     f.Filename,
		OriginalName: f.OriginalName,
		Mimetype:     f.Mimetype,
		Url:          f.Url,
		Size:         f.Size,
		Extension:    f.Extension,
		Base: Base{
			CreatedAt: f.CreatedAt,
			UpdatedAt: f.UpdatedAt,
		},
	}
}
