package lib

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

func FileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

//生成目录并拷贝文件
func FileCopyAllDir(srcPath, destPath string) (err error) {
	//检测目录正确性
	if srcInfo, err := os.Stat(srcPath); err != nil {
		return err
	} else {
		if !srcInfo.IsDir() {
			e := errors.New("srcPath不是一个正确的目录！")
			return e
		}
	}

	if destInfo, err := os.Stat(destPath); err != nil {
		err :=os.MkdirAll(destPath, os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		if !destInfo.IsDir() {
			e := errors.New("destInfo不是一个正确的目录！")
			return e
		}
	}

	err = filepath.Walk(srcPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if !f.IsDir() {
			newFile := destPath + string(os.PathSeparator) + filepath.Base(path)
			FileCopy(path, newFile)
		}
		return nil
	})

	return
}


/**
	old : /root/t.txt
	new : /dev/new.txt
 */
func FileCopy(old, new string) (err error) {
	oldFile, err := os.Open(old)
	if err != nil {
		return
	}
	defer oldFile.Close()

	newFile, err := os.Create(new)
	if err != nil {
		return
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, oldFile)
	if err != nil {
		return
	}
	err = newFile.Sync()
	return
}
