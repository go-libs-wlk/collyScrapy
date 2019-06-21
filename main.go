package main

import (
	"collyScrapy/logic"
	"collyScrapy/models"
	_ "collyScrapy/routers"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	db_host := beego.AppConfig.String("db_host")
	db_user := beego.AppConfig.String("db_user")
	db_password := beego.AppConfig.String("db_password")
	db_port := beego.AppConfig.String("db_port")
	db_name := beego.AppConfig.String("db_name")
	mysqlConnection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", db_user, db_password, db_host, db_port, db_name)
	//注册驱动
	orm.RegisterDriver("mysql", orm.DRMySQL)
	//设置默认数据库
	orm.RegisterDataBase("default", "mysql", mysqlConnection, 30)
	orm.RegisterModel(new(models.Video))
	orm.RunSyncdb("default", false, true)
}

func main() {

	//if f, _ := beego.AppConfig.Bool("is_main"); f {
	//	go logic.Start()
	//}
	//go logic.StartDownloadVideo()

	go logic.StartTrans()
	//var video models.Video
	//video.Id = 1
	//video.VideoDir = "100"
	//video.Num = "1111111"
	//video.SetVideoStatus(models.VideoWaitTrans, "video_dir", "num")
	beego.Run()
}
