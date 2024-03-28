package main

import (
	"image"
	"io"
	"io/fs"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type MultiDirSectionHelper struct {
	sourceBaseDir string
	album         string
	destBaseDir   string
}

func (sectionHelper MultiDirSectionHelper) DestBaseDir() string {
	return sectionHelper.destBaseDir
}

func (sectionHelper MultiDirSectionHelper) SourceBaseDir() string {
	return sectionHelper.sourceBaseDir
}

func (sectionHelper MultiDirSectionHelper) ScanSection() []ISection {
	dir := os.DirFS(sectionHelper.SourceBaseDir())
	sectionList := make([]ISection, 0)

	sectionMap := make(map[string]*MultiDirSection)

	dirEntityList, _ := fs.ReadDir(dir, ".")
	// sort.Slice(dirEntityList, func(i, j int) bool {
	// 	fileInfoI, _ := dirEntityList[i].Info()
	// 	fileInfoJ, _ := dirEntityList[i].Info()
	// 	return fileInfoI.ModTime().UnixMilli() < fileInfoJ.ModTime().UnixMilli()
	// })
	for _, dirEntity := range dirEntityList {
		fileInfo, _ := dirEntity.Info()
		dirStamp := fileInfo.ModTime().UTC().Format("20060102150405")

		if !dirEntity.IsDir() {
			continue
		}
		subSection := Section{}
		subSection.sourceAlbum = sectionHelper.album
		subSection.destAlbum = sectionHelper.album
		subSection.imgList = make([]Image, 0)
		subSection.name = dirEntity.Name()
		subSection.clientStatus = "NONE"
		now := time.Now()
		stamp := now.Format("20060102150405")
		subSection.timeStamp = stamp

		imgList, _ := fs.ReadDir(dir, dirEntity.Name())

		for _, img := range imgList {
			imgName := img.Name()
			if strings.HasSuffix(imgName, ".webp") || strings.HasSuffix(imgName, ".jfif") || strings.HasSuffix(imgName, ".jpg") || strings.HasSuffix(imgName, ".jpeg") || strings.HasSuffix(imgName, ".JPG") || strings.HasSuffix(imgName, ".JPEG") || strings.HasSuffix(imgName, ".png") || strings.HasSuffix(imgName, ".PNG") {
				image := Image{}
				image.name = img.Name()
				image.binName = img.Name()
				subSection.imgList = append(subSection.imgList, image)
			}
		}

		sort.Slice(subSection.imgList, func(i, j int) bool {
			var name1 = subSection.imgList[i].name
			var name2 = subSection.imgList[j].name
			var pName1 = strings.Split(name1, ".")[0]
			var pName2 = strings.Split(name2, ".")[0]
			var index1, err1 = strconv.Atoi(pName1)
			var index2, err2 = strconv.Atoi(pName2)
			if err1 != nil || err2 != nil {
				return i < j
			}
			return index1 < index2
		})
		// parse sub-sectoin name
		pureName := strings.Split(subSection.name, "-")[0]

		existSection, exist := sectionMap[pureName]
		if !exist {
			tmpExistSection := MultiDirSection{}
			tmpExistSection.subSections = make([]Section, 0)
			tmpExistSection.name = pureName
			tmpExistSection.album = sectionHelper.album
			tmpExistSection.timeStamp = dirStamp

			sectionMap[pureName] = &tmpExistSection
			existSection = sectionMap[pureName]
		}
		existSection.subSections = append(existSection.subSections, subSection)
		sort.Slice(existSection.subSections, func(i, j int) bool {
			name1 := existSection.subSections[i].name
			name2 := existSection.subSections[j].name
			indexStr1, _ := strings.CutSuffix(strings.Split(name1, "-")[1], "_files")
			indexStr2, _ := strings.CutSuffix(strings.Split(name2, "-")[1], "_files")

			index1, atoiErr1 := strconv.Atoi(indexStr1)
			index2, atoiErr2 := strconv.Atoi(indexStr2)
			if atoiErr1 != nil || atoiErr2 != nil {
				return i < j
			}

			return index1 < index2
		})

	}
	for _, section := range sectionMap {
		sectionList = append(sectionList, section)

	}
	sort.Slice(sectionList, func(i, j int) bool {
		return sectionList[i].TimeStamp() < sectionList[j].TimeStamp()
	})

	return sectionList
}

type MultiDirSection struct {
	subSections []Section
	name        string
	album       string
	timeStamp   string
}

func (section *MultiDirSection) Status() string {
	return section.subSections[0].clientStatus
}

// SaveToDb implements ISection.
func (section *MultiDirSection) SaveToDb() {

	sectionId := insertISection(section)

	imageList := section.ImageList()

	for _, image := range imageList {
		insertImg(image, sectionId)
	}

}

// Album implements ISection.
func (section MultiDirSection) Album() string {
	return section.album
}

// Cover implements ISection.
func (section MultiDirSection) Cover() Image {
	return section.subSections[0].imgList[0]
}

// ParseSize implements ISection.
func (section *MultiDirSection) ParseSize(sectionHelper ISectionHelper) {
	for _, subSection := range section.subSections {
		totalCount := len(subSection.imgList)
		for i, imageItem := range subSection.imgList {
			targetFile := (sectionHelper.DestBaseDir() + "/" + section.album + "/" + section.name + "/" + imageItem.name)

			imgReader, _ := os.Open(targetFile)
			img, _, err := image.Decode(imgReader)
			if err != nil {
				msgChan <- BatchComment{section.name, imageItem.name + ":" + err.Error()}

				continue
			}

			x := img.Bounds().Dx()
			y := img.Bounds().Dy()
			log.Default().Printf("(%d/%d) parse %s succ, height:%d, width:%d", i, totalCount, imageItem.name, y, x)

			imageItem.height = y
			imageItem.width = x
			subSection.imgList[i] = imageItem
		}
	}
	// log.Default().Panicln(section)
}

// CpSection implements ISection.
func (section MultiDirSection) CpSection(sectionHelper ISectionHelper) {
	newDir := sectionHelper.DestBaseDir() + "/" + section.album + "/" + section.name

	os.Mkdir(newDir, 0750)

	imageNameMap := make(map[string]interface{}, 0)
	for _, subSection := range section.subSections {
		for _, image := range subSection.imgList {
			_, exist := imageNameMap[image.name]
			if exist {
				continue
			}

			imageNameMap[image.name] = image
			targetFile, _ := os.Create(sectionHelper.DestBaseDir() + "/" + section.album + "/" + section.name + "/" + image.name)
			srcFile, err := os.Open(sectionHelper.SourceBaseDir() + "/" + subSection.name + "/" + image.name)

			if err != nil {
				if os.IsNotExist(err) {
					msg := image.name + " not exist"
					msgChan <- BatchComment{subSection.name, msg}

					log.Println(err)
					continue

				}
			}

			io.Copy(targetFile, srcFile)

		}
	}

}

// ImageList implements ISection.
func (section MultiDirSection) ImageList() []Image {
	// panic("unimplemented")
	imageList := make([]Image, 0)
	imageNameMap := make(map[string]interface{}, 0)
	for _, subSection := range section.subSections {
		for _, image := range subSection.imgList {
			_, exist := imageNameMap[image.name]
			if exist {
				continue
			}

			imageNameMap[image.name] = image
			imageList = append(imageList, image)
		}
	}
	return imageList
}

// TimeStamp implements ISection.
func (section MultiDirSection) TimeStamp() string {
	return section.timeStamp
}

func (section MultiDirSection) Name() string {
	return section.name
}
