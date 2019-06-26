package models

import (
	"github.com/astaxie/beego/orm"
)

type Server struct {
	Id       int64
	Ip       string
	Port     string
	Username string
	Password string
	VideoNum int
	LimitNum int
}

func GetServer() (*Server, error)  {
	o := orm.NewOrm()
	var server Server
	sql := "SELECT T0.`id`, T0.`ip`, T0.`port`, T0.`username`, T0.`password`, T0.`video_num`, T0.`limit_num` FROM `server` T0 WHERE T0.`video_num` < `limit_num` ORDER BY T0.`id` ASC LIMIT 1"
	err := o.Raw(sql).QueryRow(&server)
	return &server, err
}

func (s *Server)UpdateVideoNum()  {
	o := orm.NewOrm()
	var server Server
	server.Id = s.Id
	server.VideoNum = s.VideoNum
	_,_ = o.Update(&server, "video_num")
}


