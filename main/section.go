package main

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
}
