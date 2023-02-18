package main

import (
	"database/sql"
	"fmt"
	"log"
)

var db *sql.DB

func main() {
	initFlowDB()
	// file, err := os.Open("/mnt/2048/CLImages2/[西川康] お嬢様は戀話がお好き.html")
	fileNames := scanWebFile()

	for _, fileName := range fileNames {
		process1024Web(fileName)
		// log.Println("file", fileName)

	}

}

func process1024Web(fileName string) {
	log.Println("process1024Web", fileName)
	insertLog(fileName, "")
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
	section := cpFiles(imgSrcList, realDir)

	section = parseImage(section)

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
}

type Image struct {
	height int
	width  int
	name   string
}
