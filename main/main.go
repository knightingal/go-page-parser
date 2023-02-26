package main

import (
	"database/sql"
	"fmt"
	"log"
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
		process1024Web(fileName)
	}

}

func process1024Web(fileName string) {
	log.Println("process1024Web", fileName)

	if checkSuccLog(fileName) {
		log.Println(fileName, "succ")
		return
	}
	if !checkExistLog(fileName) {
		insertLog(fileName, "")
	}

	imgSrcList, srcDir := parseDoc(fileName)
	if len(imgSrcList) == 0 {
		updateLog(fileName, "img not found")
		return
	}

	realDir, succ := matchDirName(srcDir)

	if !succ {
		updateLog(fileName, "dir not found")
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
