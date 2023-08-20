package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/fs"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "golang.org/x/image/webp"
)

const SOURCE_DIR = "/mnt/2048/2048/CLImages/CLImages0828/"
const TARGET_DIR = "/mnt/linux1000/1805/"
const BAK_DIR = "/mnt/bak/2048/"

func generateTargetFullPath(dirName string, imgName string) string {
	return TARGET_DIR + dirName + "/" + imgName
}

func cpFiles(imgSrcList []string, realDirName string, docPath string) Section {
	now := time.Now()
	stamp := now.Format("20060102150405")

	section := Section{}
	section.timeStamp = stamp
	section.name = stamp + realDirName
	section.webName = docPath

	imgList := make([]Image, 0)

	os.Mkdir(TARGET_DIR+section.name, 0750)
	for _, imgSrc := range imgSrcList {
		targetFile, _ := os.Create(generateTargetFullPath(section.name, imgSrc))
		srcFile, err := os.Open(SOURCE_DIR + realDirName + "/" + imgSrc)
		if err != nil {
			if os.IsNotExist(err) {
				msg := imgSrc + " not exist"
				msgChan <- BatchComment{docPath, msg}

				log.Println(err)
				continue

			}
		}
		io.Copy(targetFile, srcFile)
		image := Image{}
		image.name = imgSrc
		imgList = append(imgList, image)
		srcFile.Close()
		targetFile.Close()
	}
	section.imgList = imgList

	return section
}

func scanWebFile() []string {
	dir := os.DirFS(SOURCE_DIR)
	dirEntityList, _ := fs.ReadDir(dir, ".")
	fileNames := make([]string, 0)
	for _, file := range dirEntityList {
		if strings.HasSuffix(file.Name(), ".html") && strings.HasPrefix(file.Name(), "[修复]") {
			fileNames = append(fileNames, file.Name())
		}
	}

	return fileNames

}

// match exist dir with a given dir name parse from html file
func matchDirName(srcDir string) (matchDirName string, succ bool) {
	dir := os.DirFS(SOURCE_DIR)

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
	file, err := os.Open(SOURCE_DIR + docPath)
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
			if len(srcDirList) < 2 {
				return
			}
			srcDir = srcDirList[len(srcDirList)-2]
			imgName := srcDirList[len(srcDirList)-1]
			imgSrcList = append(imgSrcList, imgName)
		})
	})
	return
}

func parseDocV2(docPath string) (imgSrcList []string, srcDir string) {
	file, err := os.Open(SOURCE_DIR + docPath)
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
			src, _ := s.Attr("ess-data")
			escape, _ := url.QueryUnescape(src)

			fmt.Println(escape)
			srcDirList := strings.Split(escape, "/")
			if len(srcDirList) < 2 {
				return
			}
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

func parseImage(section Section) (Section, error) {

	for i, imgSt := range section.imgList {
		imgReader, _ := os.Open(generateTargetFullPath(section.name, imgSt.name))
		img, _, err := image.Decode(imgReader)
		if err != nil {
			msgChan <- BatchComment{section.webName, imgSt.name + ":" + err.Error()}

			continue
		}

		x := img.Bounds().Dx()
		y := img.Bounds().Dy()

		imgSt.height = y
		imgSt.width = x

		section.imgList[i] = imgSt

	}
	section.album = "1806"
	section.cover = section.imgList[0]

	return section, nil

}
