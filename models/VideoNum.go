package models

import "github.com/astaxie/beego/orm"

type VideoNum struct {
	Id int64
	VideoNum string `orm:"unique"`
}

func VideoNumAdd(videoNum string) (id int64, err error) {
	o := orm.NewOrm()
	var vNumModel VideoNum
	vNumModel.VideoNum = videoNum
	return o.Insert(&vNumModel)
}

func VideoNumDel(videoNum string)  (id int64, err error) {
	o := orm.NewOrm()
	var v VideoNum
	v.VideoNum = videoNum
	return o.Delete(&v, "video_num")
}

func VideoNumFind(videoNum string) error {
	o := orm.NewOrm()
	var vNum VideoNum
	vNum.VideoNum = videoNum
	return o.Read(&vNum, "video_num")
}