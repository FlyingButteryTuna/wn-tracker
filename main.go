package main

import (
	"fmt"
	"github.com/FlyingButterTuna/wn-tracker/json"
	"github.com/FlyingButterTuna/wn-tracker/parsers"
	"net/http"
)

func main() {
	//url := "https://ncode.syosetu.com/n9669bk/"
	//url := "https://ncode.syosetu.com/n6524iw/"
	url := "https://ncode.syosetu.com/n8769iq/"
	client := &http.Client{}
	doc, err := parsers.FetchPage(url, client)
	if err != nil {
		fmt.Println("Errrrror")
		return
	}

	var parser parsers.NovelParser = &parsers.NarouParser{Link: url}
	title := parser.ParseTitle(doc)
	chaptersData := parser.ParseTOC(doc)
	novelData := parsers.NovelData{Title: title, Sections: chaptersData, Link: url}

	json.SaveNovelDataToFile(novelData, "test.json")
	fmt.Println(title)
}
