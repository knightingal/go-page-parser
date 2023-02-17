package main

import (
	"database/sql"
	"fmt"
)

var db *sql.DB

func main() {
	// file, err := os.Open("/mnt/2048/CLImages2/[西川康] お嬢様は戀話がお好き.html")
	imgSrcList, srcDir := parseDoc("輝夜姬想讓人告白_天才們的戀愛頭腦戰_ 早坂愛 2.html")

	realDir, succ := matchDirName(srcDir)

	if !succ {
		return
	}
	fmt.Println(realDir)
	section := cpFiles(imgSrcList, realDir)

	section = parseImage(section)

	initFlowDB()
	sectoinId := insertSection(section)
	for _, imgSt := range section.imgList {
		insertImg(imgSt, sectoinId)
	}

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
