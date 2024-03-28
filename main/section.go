package main

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"
)

type ISection interface {
	Name() string

	ImageList() []Image

	Cover() Image

	Album() string

	TimeStamp() string

	CpSection(sectionHelper ISectionHelper)

	ParseSize(sectionHelper ISectionHelper)

	SaveToDb()

	Status() string
}

type ISectionHelper interface {
	ScanSection() []ISection

	SourceBaseDir() string

	DestBaseDir() string
}
