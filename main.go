package main

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	//targz.Compress("backup", "file.tar.gz")
	//targz.Extract("file.tar.gz", "bk")
	//err := Tar("backup", ".")
	//if err != nil {
	//	fmt.Println(err)
	//}

	err := Untar("backup.tar", "bk")
	if err != nil {
		fmt.Println(err)
	}

}

func Tar(source, target string) error {
	filename := filepath.Base(source)
	target = filepath.Join(target, fmt.Sprintf("%s.tar", filename))
	tarfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			if baseDir != "" {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			}

			if err := tarball.WriteHeader(header); err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tarball, file)
			return err
		})
}

func Untar(tarball, target string) error {
	reader, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		fmt.Println(header.Uid, header.Gid)
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}
	return nil
}

//func handler(signal os.Signal) {
//	if signal == syscall.SIGTERM {
//		fmt.Println("Got kill signal. ")
//		fmt.Println("Program will terminate now.")
//		return
//	} else if signal == syscall.SIGINT {
//		fmt.Println("Got CTRL+C signal")
//		fmt.Println("Closing.")
//		return
//	} else {
//		fmt.Println("Ignoring signal: ", signal)
//	}
//}
//
//func tes() chan int {
//	sigchnl := make(chan os.Signal, 1)
//	signal.Notify(sigchnl)
//	exitchnl := make(chan int)
//	go func() {
//		for {
//			select {
//			case <-sigchnl:
//				s := <-sigchnl
//				handler(s)
//				return
//			default:
//				fmt.Println("AA")
//			}
//
//		}
//	}()
//	return exitchnl
//}
//
//func main() {
//	sigchnl := make(chan os.Signal, 1)
//	signal.Notify(sigchnl)
//	exitcode := <-tes()
//	fmt.Println("Ignoring signal: ")
//	defer func() {
//		fmt.Println("halo")
//	}()
//	<-sigchnl
//	os.Exit(exitcode)
//}
