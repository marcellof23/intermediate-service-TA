package consumer

import (
	"context"
	"io"
	"os"

	"github.com/intermediate-service-ta/boot"
	"github.com/intermediate-service-ta/helper"
)

func BackupFiletoDisk(ctx context.Context, msg Message) error {
	filepath := helper.JoinPath(boot.Backup, msg.AbsPathDest, msg.AbsPathSource)
	var osFile, err = os.Create(filepath)
	err = os.Chmod(filepath, os.FileMode(msg.FileMode))
	if err != nil {
		return err
	}

	err = os.Chown(filepath, msg.Uid, msg.Gid)
	if err != nil {
		return err
	}

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

func RemoveFileFromDisk(ctx context.Context, msg Message) error {
	filepath := helper.JoinPath(boot.Backup, msg.AbsPathSource)
	err := os.RemoveAll(filepath)
	if err != nil {
		return err
	}
	return nil
}

func RemoveFolderFromDisk(ctx context.Context, msg Message) error {
	filepath := helper.JoinPath(boot.Backup, msg.AbsPathSource)
	err := os.RemoveAll(filepath)
	if err != nil {
		return err
	}
	return nil
}

func CopyFiletoDisk(ctx context.Context, msg Message) error {
	filepathSrc := helper.JoinPath(boot.Backup, msg.AbsPathSource)
	filepathDest := helper.JoinPath(boot.Backup, msg.AbsPathDest)
	//os.MkdirAll(filepathSrc, os.FileMode(msg.FileMode))

	originalFile, err := os.Open(filepathSrc)
	if err != nil {
		return err
	}
	defer originalFile.Close()

	newFile, err := os.Create(filepathDest)
	os.Chmod(filepathDest, os.FileMode(msg.FileMode))
	os.Chown(filepathDest, msg.Uid, msg.Gid)
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
