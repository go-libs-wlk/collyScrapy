package lib

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func DirCreate(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func DirCheck(path string) bool {
	_ , err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

func DirCheckOrCreate(path string) error {
	if DirCheck(path) {
		return nil
	}
	return DirCreate(path)
}


func DirCreateByDirNumberName(path string) (dir string,err error) {
	// 创建二级目录 10个
	s, _ := ioutil.ReadDir(path)
	// 获取最大的序号
	max := 0
	for _, v := range s {
		dirName, err := strconv.Atoi(v.Name())
		if err == nil{
			if max < dirName {
				max = dirName
			}
		}
	}
	dir = path + string(os.PathSeparator) + strconv.Itoa(max+1)
	err = DirCheckOrCreate(dir)
	return
}

func DirPlusNumberNameOne(path string) (newPath string, err error) {
	lastNumberName := filepath.Base(path)
	lastNumberNameInt, err := strconv.Atoi(lastNumberName)
	if err != nil {
		return
	}
	newPath = filepath.Dir(path) + string(os.PathSeparator) + strconv.Itoa(lastNumberNameInt + 1)
	err = DirCheckOrCreate(newPath)
	return
}

func DirGetLastDirByNumberName(path string) string {
	s, err := ioutil.ReadDir(path)
	if err != nil {
		DirCheckOrCreate(path)
	}
	// 获取最大的序号
	max := 0
	for _, v := range s {
		dirName, err := strconv.Atoi(v.Name())
		if err == nil{
			if max < dirName {
				max = dirName
			}
		}
	}
	if max == 0 {
		max = 1
	}

	return path + string(os.PathSeparator) + strconv.Itoa(max)
}

func DirGetFileNumber(path string) (num int, err error) {
	s, err := ioutil.ReadDir(path)
	if err != nil{
		return
	}

	for _, v := range s {
		if !v.IsDir() {
			num++
		}
	}
	return
}

func DirGetDirNumber(path string) (num int, err error) {
	s, err := ioutil.ReadDir(path)
	if err != nil{
		return
	}

	for _, v := range s {
		if v.IsDir() {
			num++
		}
	}
	return
}