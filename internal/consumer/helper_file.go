package consumer

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/intermediate-service-ta/boot"
)

func BackupFiletoDisk(ctx context.Context, msg Message) error {
	var osFile, err = os.Create(filepath.Join(boot.Backup, filepath.Join(msg.AbsPathDest, msg.AbsPathSource)))
	if err != nil {
		return err
	}
	defer osFile.Close()

	_, err = osFile.Write(msg.Buffer)
	if err != nil {
		return err
	}

	return nil
}

func CopyFiletoDisk(ctx context.Context, pathSource, pathDest string) error {
	originalFile, err := os.Open(pathSource)
	if err != nil {
		return err
	}
	defer originalFile.Close()

	newFile, err := os.Create(pathDest)
	if err != nil {
		return err
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, originalFile)
	if err != nil {
		return err
	}

	return nil
}
