package main

import (
	"io/fs"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ISection interface {
	Name() string

	ImageList() []Image

	Cover() Image

	Album() string

	TimeStamp() string

	CpSection()
}

type ISectionHelper interface {
	ScanSection() []ISection

	SourceBaseDir() string

	DestBaseDir() string
}

type MultiDirSectionHelper struct {
	sourceBaseDir string
	album         string
}

func (sectionHelper MultiDirSectionHelper) SourceBaseDir() string {
	return sectionHelper.sourceBaseDir
}

func (sectionHelper MultiDirSectionHelper) ScanSection() []ISection {
	dir := os.DirFS(sectionHelper.SourceBaseDir())
	sectionList := make([]ISection, 0)

	sectionMap := make(map[string]*MultiDirSection)

	dirEntityList, _ := fs.ReadDir(dir, ".")
	for _, dirEntity := range dirEntityList {
		if !dirEntity.IsDir() {
			continue
		}
		subSection := Section{}
		subSection.sourceAlbum = sectionHelper.album
		subSection.destAlbum = sectionHelper.album
		subSection.imgList = make([]Image, 0)
		subSection.name = dirEntity.Name()
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

	return sectionList
}

type MultiDirSection struct {
	subSections []Section
	name        string
	album       string
}

// Album implements ISection.
func (section MultiDirSection) Album() string {
	return section.album
}

// Cover implements ISection.
func (section MultiDirSection) Cover() Image {
	return section.subSections[0].cover
}

// CpSection implements ISection.
func (section MultiDirSection) CpSection() {
}

// ImageList implements ISection.
func (section MultiDirSection) ImageList() []Image {
	// panic("unimplemented")
	imageList := make([]Image, 0)
	for _, subSection := range section.subSections {
		imageList = append(imageList, subSection.imgList...)
	}
	return imageList
}

// TimeStamp implements ISection.
func (section MultiDirSection) TimeStamp() string {
	return section.subSections[0].timeStamp
}

func (section MultiDirSection) Name() string {
	return section.name
}
