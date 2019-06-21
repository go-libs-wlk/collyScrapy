package controllers

import (
	"collyScrapy/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	url := c.GetString("url")
	var video models.Video
	video.Iframe = url
	o := orm.NewOrm()

	err := o.Read(&video, "iframe")
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}
	_, err = video.SetVideoStatus(models.VideoOk)
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	c.Data["json"] = "执行成功"
	c.ServeJSON()
}
