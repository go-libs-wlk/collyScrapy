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
	orm.RegisterModel(new(models.Video), new(models.Server))
	orm.RunSyncdb("default", false, true)
	//orm.Debug = true
}

func main() {

	if f, _ := beego.AppConfig.Bool("is_main"); f {
		go logic.Start()
	}

	go logic.StartDownloadVideo()

	go logic.StartTrans()

	beego.Run()
}
