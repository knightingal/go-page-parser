package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var db *sql.DB

var msgChan chan BatchComment

func batchCommentListener() {
	for {
		comment := <-msgChan
		updateComment(comment.fileName, comment.comment)
	}
}

func main() {

	msgChan = make(chan BatchComment)

	initFlowDB()

	go batchCommentListener()
	fileNames := scanWebFile()

	for _, fileName := range fileNames {
		process1024Web(fileName, persistenceDir)
	}
	// process1024Web("[修复][動漫] [Bird Forest (梟森)] 1RTで仲悪いノンケ女子たちが1秒キスするシリーズ-付き合ってください! [中国翻訳][34P] - 新時代的我們 草榴社區 - cl.tfrw.xyz.html", persistenceDir)

}

func bakDir(realDir, fileName string) {
	os.Rename(SOURCE_DIR+realDir, BAK_DIR+realDir)
	os.Rename(SOURCE_DIR+fileName, BAK_DIR+fileName)
}

func persistenceDir(realDir, fileName string) {
	imgSrcList, _ := parseDoc(fileName)
	if len(imgSrcList) == 0 {
		updateLog(fileName, "img not found")
		return
	}
	fmt.Println(realDir)
	section := cpFiles(imgSrcList, realDir, fileName)

	section, err := parseImage(section)
	if err != nil {
		updateLog(fileName, err.Error())
	}

	sectoinId := insertSection(section)
	for _, imgSt := range section.imgList {
		insertImg(imgSt, sectoinId)
	}
	updateLog(fileName, "succ")
}

func process1024Web(fileName string, dirProcessor func(string, string)) {
	log.Println("process1024Web", fileName)

	if checkSuccLog(fileName) {
		log.Println(fileName, "succ")
		return
	}
	if !checkExistLog(fileName) {
		insertLog(fileName, "")
	}

	_, srcDir := parseDoc(fileName)

	realDir, succ := matchDirName(srcDir)

	if !succ {
		updateLog(fileName, "dir not found")
		return
	}

	dirProcessor(realDir, fileName)
}

type Section struct {
	timeStamp string
	name      string
	album     string
	imgList   []Image
	cover     Image
	webName   string
}

type Image struct {
	height int
	width  int
	name   string
}

type BatchComment struct {
	fileName string
	comment  string
}
