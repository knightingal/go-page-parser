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
	legacyOrder := false
	test := false
	flow1000 := false
	encrytpe := false
	testMulti := true
	msgChan = make(chan BatchComment)
	initFlowDB()
	go batchCommentListener()
	if testMulti {
		multiHelper := MultiDirSectionHelper{"/home/knightingal/Downloads/20240310/", "1807"}
		sectionList := multiHelper.ScanSection()
		fmt.Println(sectionList[0].Name())
		fmt.Println(sectionList[0].ImageList())

		return
	}

	if flow1000 {
		sectionList := scanFLow1000Dir(encrytpe)
		for _, section := range sectionList {
			if len(section.imgList) == 0 {
				log.Default().Printf("[%s] scan failed", section.name)
				continue
			}
			section, _ = parseImage(section)
			sectoinId := insertSection(section)
			for _, imgSt := range section.imgList {
				insertImgBin(imgSt, sectoinId)
			}
		}
		return
	}

	if legacyOrder {

		sectionList := scanLegacyDir()
		sectionList = cpSections(sectionList)
		println(sectionList)

		return
	}

	if test {
		processFlow1000Web(
			"target.html",
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
	section := cpFiles(imgSrcList, realDir, fileName, true)

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
	sourceAlbum  string
	destAlbum    string
	imgList      []Image
	cover        Image
	webName      string
	clientStatus string
}

type Image struct {
	height       int
	width        int
	name         string
	binName      string
	milliseconds int64
}

type BatchComment struct {
	fileName string
	comment  string
}
