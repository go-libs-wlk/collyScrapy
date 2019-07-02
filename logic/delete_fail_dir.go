package logic

import (
	"collyScrapy/models"
	"fmt"
	"github.com/astaxie/beego"
	"path/filepath"
)

func DeleteDir()  {
	server := beego.AppConfig.String("server")

	videos , err := models.GetAllVideoWait(server)
	if err != nil {
		beego.Error(err.Error())
		return
	}

	for _, v := range videos {
		fmt.Printf("%+v \n", v)
		if v.VideoDir != "" {
			path, err := getVideoFile(v)
			if err != nil {
				beego.Error(err.Error())
				continue
			}
			beego.Info("删除文件:",path, " 所在的目录：", filepath.Dir(path))
			//return
			//os.RemoveAll(filepath.Dir(path))
		}
	}

}
