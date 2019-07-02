package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

const (
	VideoWait = iota + 1
	VideoDownloading
	VideoWaitTrans
	VideoTransing
	VideoDownloadFail
	VideoTransFail
	VideoUploading
	VideoUploadFail
	VideoOk
	VideoFileNotFound
)

type Video struct {
	Id        int64
	Iframe    string `orm:"unique"`
	Num       string // 车牌号
	Status    uint8  // 1 待下载 2 已下载 3 待转吗 4 已转码 5 下载失败 6 转码失败
	ServerNum int    // 服务器编号
	Label string
	ErrMsg string
	VideoDir string
	CreatedAt int64
}

func VideoAddPageUrl(url string) (id int64, err error)  {
	o := orm.NewOrm()
	var video Video
	video.Status = VideoWait
	video.Iframe = url
	video.Num = ""
	video.ServerNum = 0
	video.CreatedAt = time.Now().Unix()
	id, err = o.Insert(&video)
	return
}

// 获取待下载的视频信息
func GetVideoWaitDownload() (*Video, bool)  {
	server ,_ := beego.AppConfig.Int("server")
	videoDefaultNum := beego.AppConfig.DefaultInt64("video_num", 0)
	var video Video
	o := orm.NewOrm()
	num, _ := o.QueryTable(new(Video)).Filter("server_num", server).Count()
	if num >= videoDefaultNum {
		beego.Info("视频已经超出，database_video_num=" , num, ". 默认设置 video_num = ", videoDefaultNum)
		return nil, true
	}

	_ = o.QueryTable(new(Video)).Filter("status", VideoWait).OrderBy("-id").One(&video)
	if video.Id == 0 {
		beego.Info("暂时无可下载的视频 休息十分钟")
		time.Sleep(10 * time.Minute)
		return GetVideoWaitDownload()
	}

	video.ServerNum = server
	video.Status = VideoDownloading
	o.Update(&video, "server_num", "status")
	return &video, false
}

func (v *Video)SetVideoStatus(status uint8, cols ... string) (int64, error) {
	o := orm.NewOrm()
	v.Status = status
	cols = append(cols, "status")
	return o.Update(v, cols...)
}

func GetVideoWaitTrans() *Video {
	var video Video
	server ,_ := beego.AppConfig.Int("server")
	o := orm.NewOrm()
	_ = o.QueryTable(new(Video)).Filter("status", VideoWaitTrans).Filter("server_num",server).OrderBy("-id").One(&video)
	if video.Id == 0 {
		beego.Info("暂时无可转码的视频 休息十分钟")
		time.Sleep(10 * time.Minute)
		return GetVideoWaitTrans()
	}
	video.Status = VideoTransing
	o.Update(&video, "status")
	return  &video
}

func GetVideoById(id int64) *Video {
	var video  Video
	o := orm.NewOrm()
	video.Id = id
	o.Read(&video)
	return &video
}


func GetVideoUploadFail() *Video {
	var video Video
	server ,_ := beego.AppConfig.Int("server")
	o := orm.NewOrm()
	_ = o.QueryTable(new(Video)).Filter("status", VideoUploadFail).Filter("server_num",server).OrderBy("-id").One(&video)
	if video.Id == 0 {
		beego.Info("暂时无上传失败的视频 休息十分钟")
		time.Sleep(10 * time.Minute)
		return GetVideoWaitTrans()
	}
	return  &video
}

func GetAllVideoWait(serverId string) ([]*Video, error) {
	o := orm.NewOrm()
	var videos []*Video
	_, err := o.QueryTable(new(Video)).Filter("status", VideoWait).Filter("server", serverId).All(&videos)
	if err != nil {
		return videos, err
	}
	return videos, nil
}
