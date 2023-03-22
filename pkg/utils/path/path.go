package path

import (
	"archive/tar"
	"compress/gzip"
	"embed"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func IsPathExist(path string) (res bool) {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func MkdirPath(path string) (err error) {
	err = os.MkdirAll(path, os.ModePerm)
	return err
}

func RemovePath(path string) {
	os.RemoveAll(path)
}

func CopyFSFile(binFS embed.FS, FSFilePath, localPath string) (err error) {
	binFsFileOpen, err := binFS.Open(FSFilePath)
	if err != nil {
		return err
	}
	defer func(binFsFileOpen fs.File) {
		err := binFsFileOpen.Close()
		if err != nil {
			fmt.Printf("binFsFileOpen close error: %+v\n", err)
		}
	}(binFsFileOpen)
	binFsFileOut, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer func(binFsFileOut *os.File) {
		err := binFsFileOut.Close()
		if err != nil {
			fmt.Printf("binFsFileOut close error: %+v\n", err)
		}
	}(binFsFileOut)
	_, err = io.Copy(binFsFileOut, binFsFileOpen)
	return err
}

func TarPath(src, dst string) (err error) {
	fw, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fw.Close()

	gw := gzip.NewWriter(fw)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	return filepath.Walk(src, func(fileName string, fi fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		hdr, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return err
		}

		hdr.Name = strings.TrimPrefix(fileName, "/tmp/")

		if err = tw.WriteHeader(hdr); err != nil {
			return err
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		fr, err := os.Open(fileName)
		defer fr.Close()

		if err != nil {
			return err
		}

		n, err := io.Copy(tw, fr)
		if err != nil {
			return err
		}

		logrus.Printf("tar %s, size: %d\n", fileName, n)
		return err
	})
}
