package main

import (
	"io/fs"
	"os"
	"strings"
	"time"
)

type ISection interface {
	Name() string

	ImageList() Image

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

		// parse sub-sectoin name
		pureName := strings.Split(subSection.name, "-")[0]
		// index, atoiErr := strconv.Atoi(strings.Split(subSection.name, "-")[1])
		// if atoiErr != nil {
		// 	index = 0
		// }
		existSection, exist := sectionMap[pureName]
		if !exist {
			tmpExistSection := MultiDirSection{}
			tmpExistSection.subSections = make([]Section, 0)
			tmpExistSection.name = pureName

			sectionMap[pureName] = &tmpExistSection
			existSection = sectionMap[pureName]
		}
		existSection.subSections = append(existSection.subSections, subSection)

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
func (section MultiDirSection) ImageList() Image {
	panic("unimplemented")
}

// TimeStamp implements ISection.
func (section MultiDirSection) TimeStamp() string {
	return section.subSections[0].timeStamp
}

func (section MultiDirSection) Name() string {
	return section.name
}
