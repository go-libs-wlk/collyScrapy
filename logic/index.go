package logic

import (
	"collyScrapy/lib"
	"collyScrapy/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
	"github.com/cavaliercoder/grab"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var MaxPageNum = 0

func Start() {
	domain := "https://javcl.com/"
	beego.Info("数据爬取中")
	if !Index(domain) {
		beego.Info("已经爬取结束, 休息30分钟后继续")
		time.Sleep(30 * time.Minute)
		Start()
	}
}

// 获取所有未采集的链接地址
func Index(url string) bool {
	var (
		res       *http.Response
		err       error
		urls      []string
		repeatNum = 1
	)
	beego.Info("爬取的地址 = ", url)
	res, err = http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("#tabsall .clip-link").Each(func(i int, s *goquery.Selection) {
		url, ok := s.Attr("href")
		if ok {
			urls = append(urls, url)
		}
	})

	for _, v := range urls {
		if repeatNum >= 5 {
			return false
		}
		_, err := models.VideoAddPageUrl(v)
		if err != nil && strings.Contains(err.Error(), "1062") {
			repeatNum++
		}
	}

	time.Sleep(3 * time.Second)
	// 获取下一页地址
	var baseUrl = "https://javcl.com/page/"
	if 0 < MaxPageNum {
		reg := regexp.MustCompile(`\/([\d]+)\/`)
		s := reg.FindStringSubmatch(url)
		if len(s) == 0 {
			url = baseUrl + "2/"
		} else {
			pageNum, err := strconv.Atoi(s[1])
			if err != nil {
				return false
			}
			url = baseUrl + strconv.Itoa(pageNum + 1) + "/"
		}
	} else {
		var ok bool
		url, ok = doc.Find(".pagination .active").Next().Find("a").Attr("href")
		if !ok {
			return false
		}
	}
	return Index(url)
}

func StartDownloadVideo() {
	downloadNum, err := beego.AppConfig.Int("download")
	if err != nil {
		beego.Error("配置文件错误，缺少并发下载数量的配置如 ：\n download = 5")
		return
	}

	for i := 0; i < downloadNum; i++ {
		go forDownloadVideo(1, nil, nil)
		time.Sleep(10 * time.Second)
	}
}

// 循环下载
func forDownloadVideo(num int, video *models.Video, err error) {
	if num > 3 {
		// 设置视频下载失败
		video.ErrMsg = err.Error()
		video.SetVideoStatus(models.VideoDownloadFail, "err_msg")
		time.Sleep(1 * time.Minute)
		forDownloadVideo(1, nil, nil)
	}
	if video == nil {
		var full bool
		video, full = models.GetVideoWaitDownload()
		if full {
			return
		}
	}

	beego.Error("开始下载文件：", fmt.Sprintf("%+v", video))

	iframe, videoNum, err := GetIframeUrl(video.Iframe)

	// 检测该视频是否被发布出去在香蕉视频
	err = models.VideoNumFind(videoNum)
	if err == nil {
		video.Num = videoNum
		video.SetVideoStatus(models.VideoOk, "num")
		forDownloadVideo(1, nil, nil)
	}

	if err != nil {
		beego.Error(err)
		num++
		time.Sleep(10 * time.Second)
		forDownloadVideo(num, video, err)
	}



	videoUrl, err := GetVideoRealUrl(iframe)
	if err != nil {
		beego.Error(err)
		num++
		time.Sleep(10 * time.Second)
		forDownloadVideo(num, video, err)
	}

	// 获取文件夹存放的
	dirPath, err := GetVideoSaveDir()
	if err != nil {
		beego.Error(err)
		num++
		time.Sleep(10 * time.Second)
		forDownloadVideo(num, video, err)
	}


	if err = Download(videoUrl.File, videoUrl.Label, videoNum, dirPath); err != nil {
		beego.Error(err)
		num++
		time.Sleep(10 * time.Second)
		forDownloadVideo(num, video, err)
	}

	video.VideoDir = filepath.Base(dirPath)
	video.Num = videoNum
	video.Label = videoUrl.Label
	video.SetVideoStatus(models.VideoWaitTrans, "video_dir", "num", "label")
	time.Sleep(1 * time.Minute)
	forDownloadVideo(1, nil, err)
}

func GetVideoSaveDir()  (newDir string, err error) {
	path := beego.AppConfig.String("video_path")
	lastDir := lib.DirGetLastDirByNumberName(path)
	lib.DirCheckOrCreate(lastDir)
	dirNum, err := lib.DirGetDirNumber(lastDir)
	if err != nil {
		return
	}

	defaultNum, err := beego.AppConfig.Int("dir_video_num")
	if err != nil {
		return
	}

	if dirNum >= defaultNum {
		newDir, err = lib.DirPlusNumberNameOne(lastDir)
		if err != nil {
			return
		}
	} else {
		newDir = lastDir
	}
	return
}

// 爬取页面所有的iframe 地址
func GetIframeUrl(urlString string) (src string, videoNum string, err error) {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	timeountCtx, timeoutCancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()
	defer timeoutCancel()
	// run task list
	var res = make(map[string]string)
	err = chromedp.Run(timeountCtx,
		chromedp.Navigate(urlString),
		chromedp.WaitVisible("#videoPlayer > iframe"),
		chromedp.Attributes(`#videoPlayer > iframe`, &res, chromedp.ByID),
		chromedp.Text(".title2", &videoNum, chromedp.BySearch),
	)

	if err != nil {
		return
	}

	src, ok := res["src"]
	if !ok {
		err = errors.New("没有匹配到ifram地址:" + urlString)
		return
	}
	return src, videoNum, nil

}

func GetVideoRealUrl(iframeUrl string) (videoRealUrl DataObject, err error) {

	urlInfo, err := url.Parse(iframeUrl)
	if err != nil {
		return
	}

	domain := urlInfo.Scheme + "://" + urlInfo.Host
	urlSlice := strings.Split(iframeUrl, "/")
	videoSign := urlSlice[len(urlSlice)-1]

	sendUrl := domain + "/api/source/" + videoSign
	//// 获取视频播放地址
	clineHttp := new(http.Client)
	data := make(url.Values)
	data.Set("r", "")
	data.Set("d", urlInfo.Host)
	request, err := http.NewRequest("POST", sendUrl, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("--->创建req失败：", err)
		return
	}
	request.Header.Add("referer", iframeUrl)
	request.Header.Add("origin", domain)
	request.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36")
	request.Header.Add("x-requested-with", "XMLHttpRequest")
	request.Header.Add("accept", "*/*")

	res, err := clineHttp.Do(request)
	if err != nil {
		return
	}
	defer res.Body.Close()
	var videoUrl = VideoUrl{}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		beego.Error("http.Do failed,[err=%s][url=%s]", err, res.Request.URL.String())
		return
	}

	err = json.Unmarshal(b, &videoUrl)
	if err != nil {
		beego.Error("结果解析错误：", err)
		return
	}

	// 获取720P地址
	beego.Info("获取的视频地址：", fmt.Sprintf("%+v", videoUrl))
	var p = []string{"720", "480", "360", "1080"}

	for _, v := range p {
		for i, vv := range videoUrl.Data {
			if strings.Contains(vv.Label, v) {
				videoRealUrl = videoUrl.Data[i]
				goto RET
			}
		}
	}
	RET:
	return
}

func Download(videoHref, videoLabel, videoNum, dirPath string) (err error) {
	beego.Info("开始下载 = ", videoHref)

	dirPath = dirPath + string(os.PathSeparator) + videoNum
	lib.DirCheckOrCreate(dirPath)

	fileName := fmt.Sprintf(dirPath + string(os.PathSeparator) + "%s-%s.mp4", videoNum, videoLabel)
	beego.Info("文件保存路径 = ", fileName)

	client := grab.NewClient()
	req, err := grab.NewRequest(fileName, videoHref)
	if err != nil {
		return
	}
	resp := client.Do(req)
	if err := resp.Err(); err != nil {
		beego.Error(fmt.Fprintf(os.Stderr, "Download failed: %v\n", err))
		return err
	}
	return
}
