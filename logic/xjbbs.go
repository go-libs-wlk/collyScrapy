package logic

import (
	"collyScrapy/models"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
	"log"
	"net/http"
	"strings"
	"time"
)

var XJBaseHttp = "http://www.xjbbs.live/"


func XjVideo(href string) bool {
	var (
		res       *http.Response
		err       error
		urls      []string
		repeatNum = 1
	)
	beego.Info("爬取的地址 = ", href)
	res, err = http.Get(href)
	if err != nil {
		beego.Error(err)
		time.Sleep(1 * time.Minute)
		return XjVideo(href)
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

	doc.Find(".subject").Each(func(i int, s *goquery.Selection) {
		url := s.Find("a").Text()
		rsUrlSlice := strings.Split(url, " ")
		urls = append(urls, rsUrlSlice[0])
	})

	for _, v := range urls {
		if repeatNum >= 5 {
			return false
		}
		_, err := models.VideoNumAdd(v)
		if err != nil && strings.Contains(err.Error(), "1062") {
			repeatNum++
		}
	}

	time.Sleep(1 * time.Second)
	// 获取下一页地址
	str, ok := doc.Find(".page-item.active").Next().Find("a").Attr("href")
	if !ok {
		return false
	}
	return XjVideo(XJBaseHttp + str)
}