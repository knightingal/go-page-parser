package main

type Flow1000SectionHelper struct {
	sourceBaseDir string
	album         string
	destBaseDir   string
}

func (sectionHelper Flow1000SectionHelper) DestBaseDir() string {
	return sectionHelper.destBaseDir
}

func (sectionHelper Flow1000SectionHelper) SourceBaseDir() string {
	return sectionHelper.sourceBaseDir
}

func (sectionHelper Flow1000SectionHelper) ScanSection() []ISection {
	return nil
}
