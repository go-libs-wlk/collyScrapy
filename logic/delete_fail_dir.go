package logic

import (
	"collyScrapy/models"
	"github.com/astaxie/beego"
	"os"
	"path/filepath"
)

func DeleteDir()  {
	server := beego.AppConfig.String("server")
	beego.Error("开始删除")
	videos , err := models.GetAllVideoWait(server)
	if err != nil {
		beego.Error(err.Error())
		return
	}
	beego.Error("共有", len(videos) , "个文件需要删除")
	for _, v := range videos {
		if v.VideoDir != "" {
			path, err := getVideoFile(v)
			if err != nil {
				beego.Error(err.Error())
				continue
			}
			beego.Info("删除文件:",path, " 所在的目录：", filepath.Dir(path))
			os.RemoveAll(filepath.Dir(path))
		}
	}

}

func DeleteOkVideoDir()  {
	server := beego.AppConfig.String("server")
	beego.Error("开始删除")
	videos , err := models.GetAllVideoOKAndDirNotDel(server)
	if err != nil {
		beego.Error(err.Error())
		return
	}
	beego.Error("共有", len(videos) , "个文件需要删除")
	for _, v := range videos {
		if v.VideoDir != "" {
			path, _ := getVideoFile(v)
			//if err != nil {
			//	beego.Error(err.Error())
			//	continue
			//}
			beego.Info("删除文件:",path, " 所在的目录：", filepath.Dir(path))
			os.RemoveAll(filepath.Dir(path))
			v.SetIsDelDirTrue()
		}
	}

}
