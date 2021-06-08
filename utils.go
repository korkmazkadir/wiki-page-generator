package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const wikiXMLfile = "wiki.xml"

func loadWikiXML(rootDirectory string) wiki {

	file, err := os.Open(rootDirectory + "/" + wikiXMLfile)
	if err != nil {
		// put a meaningfull error message
		panic(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		// put a meaningfull error message
		panic(err)
	}

	var wikiData wiki

	if err := xml.Unmarshal(data, &wikiData); err != nil {
		// put a meaningfull error message
		panic(err)
	}

	return wikiData
}

func getWikiFiles(rootDirectory string) []string {

	var fileList []string

	files, err := ioutil.ReadDir(rootDirectory)
	if err != nil {
		panic(err)
	}

	for _, f := range files {

		if f.IsDir() == false && strings.Contains(strings.ToLower(f.Name()), ".md") {

			if err != nil {
				continue
			}

			fileName := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
			fileList = append(fileList, fileName)
			fmt.Printf("markdown file: %s\n", fileName)
		}
	}

	return fileList
}

//------- Home Page Data Related -------//

func generateHomePageData(files []string, wikiData wiki) homePage {

	timeString := time.Now().Format(time.RFC1123)
	homePage := homePage{Time: timeString}

	homePage.CategoryList = generateCategoryList(wikiData)
	homePage.TagList = generatTagList(wikiData)
	homePage.DateList = generatDateList(wikiData)
	homePage.PagesWithoutEntry = getPageListWithoutEntry(files, wikiData)
	homePage.EntriesWithoutPage = getEntriesWithoutPage(files, wikiData)

	return homePage
}

func generateCategoryList(wikiData wiki) []category {
	var list []category

	// category name --> entries
	categoryMap := make(map[string][]entry)
	for _, e := range wikiData.Entries {
		entries := categoryMap[e.Category]
		entries = append(entries, e)
		categoryMap[e.Category] = entries
	}

	for k, v := range categoryMap {
		// sort slice values
		sort.Slice(v, func(i, j int) bool {
			return v[i].Title < v[j].Title
		})

		list = append(list, category{Name: k, Entries: v})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	return list
}

func generatTagList(wikiData wiki) []tag {
	var list []tag

	tagMap := make(map[string][]entry)
	for _, e := range wikiData.Entries {

		for _, t := range e.Tags {
			entries := tagMap[t]
			entries = append(entries, e)
			tagMap[t] = entries
		}

	}

	for k, v := range tagMap {
		// sort slice values
		sort.Slice(v, func(i, j int) bool {
			return strings.ToLower(v[i].Title) < strings.ToLower(v[j].Title)
		})

		list = append(list, tag{Name: k, Entries: v})
	}

	sort.Slice(list, func(i, j int) bool {
		return strings.ToLower(list[i].Name) < strings.ToLower(list[j].Name)
	})

	return list
}

func generatDateList(wikiData wiki) []date {
	var list []date

	dateMap := make(map[string][]entry)
	for _, e := range wikiData.Entries {
		entries := dateMap[e.Date]
		entries = append(entries, e)
		dateMap[e.Date] = entries
	}

	for k, v := range dateMap {
		// sort slice values
		sort.Slice(v, func(i, j int) bool {
			return strings.ToLower(v[i].Title) < strings.ToLower(v[j].Title)
		})

		list = append(list, date{Name: k, Entries: v})
	}

	sort.Slice(list, func(i, j int) bool {
		return strings.ToLower(list[i].Name) < strings.ToLower(list[j].Name)
	})

	return list
}

// returns index
func findInEntries(entries []entry, page string) (entry, bool) {

	for _, e := range entries {
		if e.Page == page {
			return e, true
		}
	}

	return entry{}, false
}

func findInStrings(list []string, value string) (string, bool) {

	for _, s := range list {
		if s == value {
			return s, true
		}
	}

	return "", false
}

func getPageListWithoutEntry(files []string, wikiData wiki) []string {
	var list []string

	for _, f := range files {

		_, ok := findInEntries(wikiData.Entries, f)

		if ok == false {
			list = append(list, f)
		}

	}

	return list
}

func getEntriesWithoutPage(files []string, wikiData wiki) []entry {

	var list []entry

	for _, e := range wikiData.Entries {

		_, ok := findInStrings(files, e.Page)

		if ok == false {
			list = append(list, e)
		}
	}

	return list
}

func calculateStats(files []string, wikiData wiki, pageData *homePage) {

	pageData.WikiStats.FileCount = len(files)
	pageData.WikiStats.TotalWordCount = 0
	pageData.WikiStats.DistinctWordCount = 0
	pageData.WikiStats.CategoryCount = len(pageData.CategoryList)
	pageData.WikiStats.TagCount = len(pageData.TagList)
	pageData.WikiStats.DateCount = len(pageData.DateList)
	pageData.WikiStats.PagesWithoutXML = len(pageData.PagesWithoutEntry)
	pageData.WikiStats.XMLWithoutPage = len(pageData.EntriesWithoutPage)

}

//------- Generate Home Page -------//

func generateHome(pageData homePage) {

	file, err := os.Open("./home-page.template")
	defer file.Close()

	if err != nil {
		// put a meaningfull error message
		panic(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		// put a meaningfull error message
		panic(err)
	}

	t, err := template.New("home-page").Parse(string(data))
	if err != nil {
		// put a meaningfull error message
		panic(err)
	}

	out, err := os.OpenFile("./Home.md", os.O_RDWR|os.O_CREATE, 0755)

	writer := bufio.NewWriter(out)

	err = t.Execute(writer, pageData)
	if err != nil {
		panic(err)
	}

	writer.Flush()

}
