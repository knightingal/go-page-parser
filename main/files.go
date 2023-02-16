package main

import (
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const BASE_DIR = "/mnt/download/"
const TARGET_DIR = "/mnt/linux1000/1024/"

func cpFiles(imgSrcList []string, realDirName string) {
	os.Mkdir(TARGET_DIR+realDirName, 0750)
	for _, imgSrc := range imgSrcList {
		targetFile, _ := os.Create(TARGET_DIR + realDirName + "/" + imgSrc)
		srcFile, _ := os.Open(BASE_DIR + realDirName + "/" + imgSrc)
		io.Copy(targetFile, srcFile)
	}
}

// match exist dir with a given dir name parse from html file
func matchDirName(srcDir string) (matchDirName string, succ bool) {
	dir := os.DirFS(BASE_DIR)

	dirEntityList, _ := fs.ReadDir(dir, ".")

	dirNames := make([]string, 0)
	for _, dir := range dirEntityList {
		if dir.IsDir() {
			dirNames = append(dirNames, dir.Name())
		}
	}
	fmt.Println(dirNames)

	cb := func(src string) (string, bool) {

		filterRet := filter(&dirNames, func(dirName string) bool {
			return strings.Contains(dirName, src)
		})

		if len(*filterRet) == 1 {
			fmt.Println("====matched====")
			fmt.Println((*filterRet)[0])
			return (*filterRet)[0], true
		}

		return "", false
	}
	matchDirName, succ = windowString(srcDir, cb)
	return

}

// parse doc from a given html file path
func parseDoc(docPath string) (imgSrcList []string, srcDir string) {
	file, err := os.Open(BASE_DIR + docPath)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	defer file.Close()

	imgSrcList = make([]string, 0)

	doc.Find(".tpc_content").Each(func(i int, s *goquery.Selection) {
		s.Find("img").Each(func(i int, s *goquery.Selection) {
			src, _ := s.Attr("src")
			escape, _ := url.QueryUnescape(src)

			fmt.Println(escape)
			srcDirList := strings.Split(escape, "/")
			srcDir = srcDirList[len(srcDirList)-2]
			imgName := srcDirList[len(srcDirList)-1]
			imgSrcList = append(imgSrcList, imgName)
		})
	})
	return

}

func windowString(src string, process func(string) (string, bool)) (string, bool) {
	srcArray := []rune(src)
	size := len(srcArray)
	stop := false
	for i := 0; i < size; i++ {
		for j := 0; j <= i; j++ {
			sub1 := srcArray[j : j+size-i]
			fmt.Println(string(sub1))
			realDir, stop := process(string(sub1))
			if stop {
				return realDir, true
			}
		}
		if stop {
			break
		}
	}
	return "", false
}

func filter[T any](src *[]T, fn func(T) bool) *[]T {
	ret := make([]T, 0)
	for _, item := range *src {
		if fn(item) {
			ret = append(ret, item)
		}
	}
	return &ret
}
