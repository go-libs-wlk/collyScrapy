package main

import (
	"collyScrapy/logic"
	"collyScrapy/models"
	_ "collyScrapy/routers"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/urfave/cli"
	"log"
	"os"
	"sort"
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
	orm.RegisterModel(new(models.Video), new(models.Server), new(models.VideoNum))
	orm.RunSyncdb("default", false, true)
	//orm.Debug = true
}

func main() {

	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:    "xjbbs",
			Aliases: []string{"xj"},
			Usage:   "采集香蕉视频的车牌号",
			Action:  func(c *cli.Context) error {
				beego.Info("开始采集香蕉视频")
				logic.XjVideo(logic.XJBaseHttp)
				return nil
			},
		},
		{
			Name:    "download",
			Aliases: []string{"d"},
			Usage:   "下载javcl.com视频，转码，压缩，上传ftp",
			Action:  func(c *cli.Context) error {
				beego.Info("开始下载javcl视频")
				go logic.StartDownloadVideo()
				go logic.StartTrans()
				select {}
			},
		},
		{
			Name:    "javcl",
			Aliases: []string{"jav"},
			Usage:   "采集javcl.com站列表页链接，存入数据库",
			Action:  func(c *cli.Context) error {
				beego.Info("开始采集javcl页面链接")
				logic.Start()
				return nil
			},
		},
		{
			Name:    "main",
			Aliases: []string{"m"},
			Usage:   "采集javcl.com站列表页链接，存入数据库,下载视频，转码，压缩，上传ftp",
			Action:  func(c *cli.Context) error {
				beego.Info("开始采集javcl.com站列表页链接，存入数据库,下载视频，转码，压缩，上传ftp")
				go logic.Start()
				go logic.StartDownloadVideo()
				go logic.StartTrans()
				select {}
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
