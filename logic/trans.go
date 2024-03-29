package logic

import (
	"bytes"
	"collyScrapy/lib"
	"collyScrapy/models"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func StartTrans() {

	transNum , err := beego.AppConfig.Int("transcode")
	if err != nil {
		beego.Error("缺少参数配置:transcode， 如：transcode = 5")
		return
	}

	for i := 0; i < transNum; i++ {
		go forTrans()
		time.Sleep(10 * time.Second)
	}

}

func forTrans() {
	video := models.GetVideoWaitTrans()
	//video := models.GetVideoById(48535) // 测试使用
	beego.Info("开始处理视频", fmt.Sprintf("%+v", video))

	// 获取video的文件地址
	videoFile, err := getVideoFile(video)
	if err != nil {
		video.ErrMsg = err.Error()
		video.SetVideoStatus(models.VideoTransFail, "err_msg")
		forTrans()
	}

	// 复制文件
	src := beego.AppConfig.String("video_ad")
	dest := filepath.Dir(videoFile)
	err = lib.FileCopyAllDir(src, dest)
	if err != nil {
		video.ErrMsg = err.Error()
		video.SetVideoStatus(models.VideoTransFail, "err_msg")
		forTrans()
	}

	// 转码
	err, stdout, stderr := transcode(videoFile)
	if err != nil {
		beego.Error(err,stdout,stderr)
		video.ErrMsg = "转码加水印" + err.Error() + "\n" + stdout + "\n" + stderr
		video.SetVideoStatus(models.VideoTransFail, "err_msg")
		forTrans()
	}

	// 拼接广告头
	err, stdout, stderr, outFile := conactAdVideo(videoFile)
	if err != nil {
		beego.Error(err,stdout,stderr)
		video.ErrMsg = "转码加水印" + err.Error() + "\n" + stdout + "\n" + stderr
		video.SetVideoStatus(models.VideoTransFail, "err_msg")
		forTrans()
	}
	video.SetVideoStatus(models.VideoUploading)
	// 上传至种源
	server, err := models.GetServer()
	if err != nil {
		video.ErrMsg = "找不到种源服务器：" + err.Error()
		video.SetVideoStatus(models.VideoUploadFail, "err_msg")
		return
	}

	// 清楚垃圾文件
	os.Remove(getWaiterVideoFile(videoFile))
	os.Remove(videoFile)
	os.Remove(filepath.Dir(videoFile) + string(os.PathSeparator) + "files.txt")
	os.Remove(filepath.Dir(videoFile) + string(os.PathSeparator) + "ad.mp4")

	var filesNeedUpload []string
	err = filepath.Walk(filepath.Dir(outFile), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			filesNeedUpload = append(filesNeedUpload, path)
		}
		return nil
	})

	err = UploadFile(server, filesNeedUpload, video.Num)

	if err == nil {
		os.RemoveAll(filepath.Dir(videoFile))
		video.SetVideoStatus(models.VideoOk)
	} else {
		video.ErrMsg = err.Error()
		video.SetVideoStatus(models.VideoUploadFail, "err_msg")
	}
	forTrans()
}

func conactAdVideo(videoFile string)  (err error, outStr, outErr, outFile string) {
	adVideoFile := beego.AppConfig.String("video_ad") + string(os.PathSeparator) + "ad.mp4"
	waterVideoFile := getWaiterVideoFile(videoFile)
	dir := filepath.Dir(videoFile)
	domain := beego.AppConfig.String("domain")
	var file *os.File

	filePathTxt := dir + string(os.PathSeparator) + "files.txt"
	file, err = os.OpenFile(filePathTxt, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		beego.Error("创建写入拼接文件错误", err.Error())
		return
	}
	defer file.Close()

	file.WriteString("file '" + strings.TrimSpace(adVideoFile) + "'\n")
	file.WriteString("file '" + strings.TrimSpace(waterVideoFile) + "'\n")
	outFile = dir + string(os.PathSeparator) + domain + filepath.Base(videoFile)
	command := "ffmpeg -y -f concat -safe 0 -i " + filePathTxt + " -c copy " + outFile
	err, outStr, outErr = ExecShell(command)
	return
}

func transcode(srcvideo string) (error, string, string) {
	logo := beego.AppConfig.String("logo")
	outFile := getWaiterVideoFile(srcvideo)
	commad := "ffmpeg -y -i " + srcvideo +" -movflags +faststart -r 25 -g 50 -crf 28 -me_method hex -trellis 0 -bf 8 -acodec aac -strict -2 -ar 44100 -ab 128k -vf \"movie=" + logo + "[watermark];[in][watermark]overlay=main_w-overlay_w-10:10[out]\" -s 1280:720 " + outFile
	return ExecShell(commad)
}


func getWaiterVideoFile(videoFile string) string {
	return filepath.Dir(videoFile) + string(os.PathSeparator) + "watermark.mp4"
}


func getVideoFile(video *models.Video) (videoFile string, err error){
	path := beego.AppConfig.String("video_path")
	videoFile = path + string(os.PathSeparator) + video.VideoDir +
		string(os.PathSeparator) + video.Num + string(os.PathSeparator) + video.Num + "-" + video.Label + ".mp4"

	if ok := lib.FileExist(videoFile); !ok {
		err = errors.New("videoFile is not found" + ":" + videoFile )
	}
	return
}

// 执行Shell 命令
func ExecShell(command string) (err error, out, errMsg string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		command = "/C " + command
		args := strings.Split(command, " ")
		cmd = exec.Command("cmd", args...)
	} else {
		cmd = exec.Command("bash", "-c", command)
	}
	// 输出转码命令
	beego.Info("执行命令：", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	out = stdout.String()
	errMsg = stderr.String()
	err = cmd.Run()
	return
}