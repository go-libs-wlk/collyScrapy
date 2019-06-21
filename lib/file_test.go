package lib

import (
	"fmt"
	"testing"
)

func TestFileCopy(t *testing.T) {

	old := "/Users/like/coding/go/src/collyScrapy/ad/【左侧打勾】香蕉社区最新地址.html"
	newfile :=  "/Users/like/coding/go/src/collyScrapy/ad/test.html"

	err := FileCopy(old, newfile)
	if err != nil {
		fmt.Println(err)
	}

}

func TestFileCopyAllDir(t *testing.T) {
	old := "/Users/like/coding/go/src/collyScrapy/ad"
	newfile :=  "/Users/like/coding/go/src/collyScrapy/ads"

	err := FileCopyAllDir(old, newfile)
	if err != nil {
		fmt.Println(err)
	}
}
