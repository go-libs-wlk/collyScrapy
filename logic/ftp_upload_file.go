package logic

import (
	"collyScrapy/models"
	"fmt"
	"github.com/dutchcoders/goftp"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Fpt struct {
	Ip       string
	Port     string
	UserName string
	Password string
}

var (
	baseDir string = "/videos"
)

func UploadFile(server *models.Server, srcPath , videoNum string) (err error) {
	var (
		ftpHandle *goftp.FTP
		path string
	)
	if ftpHandle, err = goftp.Connect(fmt.Sprintf("%s:%s", server.Ip, server.Port)); err != nil {
		return
	}
	defer ftpHandle.Quit()

	if err = ftpHandle.Login(server.Username, server.Password); err != nil {
		return
	}
	// 确保根目录存在
	lines, err := ftpHandle.List(baseDir)
	if err != nil {
		ftpHandle.Mkd(baseDir + "/1")
	} else {
		path = GetLastDirByDirNumberName(lines)
	}

	num, err := GetPathFileOrDirNum(ftpHandle, path, "dir")
	if err != nil {
		_ = ftpHandle.Mkd(path)
	}

	if num >= 1 {
		lastName , _ := strconv.Atoi(filepath.Base(path))
		path = baseDir + "/" + strconv.Itoa(lastName + 1)
		_ = ftpHandle.Mkd(path)
	}

	var file *os.File
	if file, err = os.Open(srcPath); err != nil {
		return
	}
	videoBasePath := path + "/" + videoNum
	ftpHandle.Mkd(videoBasePath)
	if err = ftpHandle.Stor(videoBasePath + "/" + filepath.Base(srcPath),file); err != nil {
		return
	}
	defer file.Close()

	server.VideoNum = server.VideoNum + 1
	server.UpdateVideoNum()
	return
}


func GetPathFileOrDirNum(ftp *goftp.FTP, path ,fileType string) (num int, err error){
	lines, err := ftp.List(path)
	if err != nil {
		if err = ftp.Mkd(path); err != nil {
			return
		}
	}
	for _, line := range lines {
		file := &FtpFile{}
		file.parseLine(line)
		if file.Type == fileType {
			num++
		}
	}
	return
}

func GetLastDirByDirNumberName(lines []string) (path string) {
	// 获取最大的序号
	file := &FtpFile{}
	max := 0
	for _, v := range lines {
		file.parseLine(v)
		if file.Type == "dir" {
			dirName, err := strconv.Atoi(file.FileName)
			if err == nil{
				if max < dirName {
					max = dirName
				}
			}
		}
	}
	if max == 0 {
		max = 1
	}
	return baseDir + "/" + strconv.Itoa(max)
}




func (f *FtpFile)parseLine(line string) {
	for _, v := range strings.Split(line, ";") {
		v2 := strings.Split(v, "=")
		switch v2[0] {
		case "perm":
			f.Perm = v2[1]
		case "type":
			f.Type = v2[1]
		case "size":
			f.Size, _ = strconv.Atoi(v2[1])
		case "modify":
			f.Modify, _ = strconv.ParseUint(v2[1], 10, 64)
		default:
			f.FileName = v[1 : len(v)-2]
		}
	}
	return
}

type FtpFile struct {
	Perm string
	Type string
	FileName string
	Size int
	Modify uint64
}