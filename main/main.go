package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	test := false
	flow1000 := false
	msgChan = make(chan BatchComment)
	initFlowDB()
	go batchCommentListener()
	if flow1000 {
		sectionList := scanFLow1000Dir()
		for _, section := range sectionList {
			section, _ = parseImage(section)
			sectoinId := insertSection(section)
			for _, imgSt := range section.imgList {
				insertImgBin(imgSt, sectoinId)
			}
		}
		return
	}

	if test {
		processFlow1000Web(
			"Lilit.A.A.Time.For.Pleasure[86P] - 新時代的我們   草榴社區 - t66y.com.html",
			persistenceDir)
	} else {
		fileNames := scanWebFile()

		for _, fileName := range fileNames {
			processFlow1000Web(fileName, persistenceDir)
		}
	}

}

func bakDir(realDir, fileName string) {
	os.Rename(SOURCE_DIR+realDir, BAK_DIR+realDir)
	os.Rename(SOURCE_DIR+fileName, BAK_DIR+fileName)
}

func persistenceDir(realDir, fileName string) {
	imgSrcList := parseDocV2(fileName, realDir)
	if len(imgSrcList) == 0 {
		updateLog(fileName, "img not found")
		return
	}
	fmt.Println(realDir)
	section := cpFiles(imgSrcList, realDir, fileName)

	section, err := parseImage(section)

	if err != nil {
		updateLog(fileName, err.Error())
		return
	}

	tx, err := db.Begin()
	if err != nil {
		updateLog(fileName, err.Error())
		return
	}

	sectoinId := insertSection(section)
	for _, imgSt := range section.imgList {
		insertImg(imgSt, sectoinId)
	}
	updateLog(fileName, "succ")
	err = tx.Commit()
	if err != nil {
		updateLog(fileName, err.Error())
	}
}

// fileName: target html file name under TARGET_DIR
func processFlow1000Web(fileName string, dirProcessor func(string, string)) {
	log.Println("process1024Web", fileName)

	if checkSuccLog(fileName) {
		log.Println(fileName, "succ")
		return
	}
	if !checkExistLog(fileName) {
		insertLog(fileName, "")
	}

	// _, srcDir := parseDocV2(fileName)
	srcDir := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	realDir, succ := matchDirName(srcDir)

	if !succ {
		updateLog(fileName, "dir not found")
		return
	}

	dirProcessor(realDir, fileName)
}

type Section struct {
	timeStamp    string
	name         string
	album        string
	imgList      []Image
	cover        Image
	webName      string
	clientStatus string
}

type Image struct {
	height  int
	width   int
	name    string
	binName string
}

type BatchComment struct {
	fileName string
	comment  string
}
