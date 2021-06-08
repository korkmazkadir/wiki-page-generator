package main

type wiki struct {
	Entries []entry `xml:"entries>entry"`
}

type entry struct {
	Date     string   `xml:"date"`
	Category string   `xml:"category"`
	Title    string   `xml:"title"`
	Link     string   `xml:"link"`
	Tags     []string `xml:"tags>string"`
	Page     string   `xml:"page"`
}

// home page related

// articles groups according to category
type category struct {
	Name    string
	Entries []entry
}

// articles grouped according to tag list
type tag struct {
	Name    string
	Entries []entry
}

// articles grouped according to year, moth and day
type date struct {
	Name    string
	Entries []entry
}

type stats struct {
	FileCount         int
	TotalWordCount    int
	DistinctWordCount int
	CategoryCount     int
	TagCount          int
	DateCount         int
	PagesWithoutXML   int
	XMLWithoutPage    int
}

// the data structure to feed the template
type homePage struct {
	// timestamp page generation
	Time               string
	CategoryList       []category
	TagList            []tag
	DateList           []date
	PagesWithoutEntry  []string
	EntriesWithoutPage []entry
	WikiStats          stats
}
