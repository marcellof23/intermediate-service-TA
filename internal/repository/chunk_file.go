package repository

import (
	"context"

	"github.com/intermediate-service-ta/internal/model"
)

type ChunkFileRepository interface {
	Create(c context.Context, chunkFile *model.ChunkFile) (model.ChunkFile, error)
	Get(c context.Context, filename string) (model.ChunkFile, error)
	GetChunkFileByFileID(c context.Context, fid int64) ([]model.ChunkFile, error)
	DeleteChunkFileByFileID(c context.Context, fid int64) error
}
