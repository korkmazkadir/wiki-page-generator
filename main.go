package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	rootDirector := flag.String("wiki-root", "./", "the path of the wiki's root director")
	flag.Parse()

	if _, err := os.Stat(*rootDirector); os.IsNotExist(err) {
		panic(fmt.Errorf("Provided directory %s does not exists", *rootDirector))
	}

	fmt.Printf("root directory is %s \n", *rootDirector)

	pages := getWikiFiles(*rootDirector)

	fmt.Printf("pages: %s\n", pages)

	wikiData := loadWikiXML(*rootDirector)

	fmt.Printf("wiki data: %s\n", wikiData)

	fmt.Println("----> Home Page Data <---æ")

	pageData := generateHomePageData(pages, wikiData)

	calculateStats(pages, wikiData, &pageData)

	generateHome(pageData)

}

/*
func generateHome(pages []wikiPage) {

	home := homePage{Time: time.Now().Format(time.RFC1123), Pages: pages}

	file, err := os.Open("./home-page.template")
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

	err = t.Execute(os.Stdout, home)
	if err != nil {
		panic(err)
	}

}
*/
