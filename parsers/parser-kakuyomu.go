package parsers

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"
)

type KakuyomuParser struct {
	Link            string
	apolloStateJson map[string]interface{}
}

const timeLayoutKakuyomu = "2006-01-02T15:04:05Z"

func (p *KakuyomuParser) ParseTitle(doc *goquery.Document) (string, error) {
	if len(p.apolloStateJson) == 0 {
		p.initializeJson(doc)
	}
	novelDataJson, ok := p.apolloStateJson[p.valueWorkId()].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("error parsing novel data json")
	}
	novelTitle, ok := novelDataJson["title"].(string)
	if !ok {
		return "", fmt.Errorf("error parsing title from json")
	}
	return novelTitle, nil
}

func (p *KakuyomuParser) ParseTOC(doc *goquery.Document) ([]SectionData, error) {
	if len(p.apolloStateJson) == 0 {
		p.initializeJson(doc)
	}

	sections := make([]SectionData, 0)

	toc, ok := p.apolloStateJson[p.valueWorkId()].(map[string]interface{})["tableOfContents"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("error parsing novel data json")
	}

	emptyTocStr := "TableOfContentsChapter:"
	var is_sectionless bool
	if firstTocEntry := toc[0].(map[string]interface{})["__ref"].(string); len(toc) == 1 && firstTocEntry == emptyTocStr {
		is_sectionless = true
	}

	for _, sectionId := range toc {
		section, _ := p.apolloStateJson[sectionId.(map[string]interface{})["__ref"].(string)].(map[string]interface{})

		chapters := make([]ChapterData, 0)
		chapterIds := section["episodeUnions"].([]interface{})
		for _, chapterId := range chapterIds {
			chapter := chapterId.(map[string]interface{})["__ref"].(string)
			chapterInfo := p.apolloStateJson[chapter].(map[string]interface{})
			chapterTitle := chapterInfo["title"].(string)
			publishedAtStr := chapterInfo["publishedAt"].(string)
			episodeId := chapterInfo["id"].(string)

			chapterLink := p.Link + "/episodes/" + episodeId
			datePusblished, _ := time.Parse(timeLayoutKakuyomu, publishedAtStr)

			chapterData := ChapterData{Name: chapterTitle, Link: chapterLink, DatePosted: datePusblished}
			chapters = append(chapters, chapterData)
		}

		var sectionTitle string
		var sectionLevel uint8
		if !is_sectionless {
			sectionInfoId := section["chapter"].(map[string]interface{})["__ref"].(string)
			sectionInfo := p.apolloStateJson[sectionInfoId].(map[string]interface{})
			sectionTitle = sectionInfo["title"].(string)
			sectionLevel = uint8(sectionInfo["level"].(float64))
		} else {
			sectionTitle = "default"
			sectionLevel = uint8(1)
		}

		sectionData := SectionData{Name: sectionTitle, Chapters: chapters, Level: sectionLevel}
		sections = append(sections, sectionData)
	}

	return sections, nil
}

func (p *KakuyomuParser) initializeJson(doc *goquery.Document) error {
	jsonDataString := doc.Find(`[type="application/json"]`).Text()

	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonDataString), &jsonData); err != nil {
		return fmt.Errorf("error parsing JSON data: %w", err)
	}

	props, ok := jsonData["props"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("error parsing 'props' from JSON data")
	}

	pageProps, ok := props["pageProps"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("error parsing 'pageProps' from JSON data")
	}

	apolloState, ok := pageProps["__APOLLO_STATE__"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("error parsing '__APOLLO_STATE__' from JSON data")
	}

	p.apolloStateJson = apolloState
	return nil
}

func (p *KakuyomuParser) valueWorkId() string {
	workId := p.Link[strings.LastIndex(p.Link[:strings.LastIndex(p.Link, "/")-1], "/")+1:]
	workId = strings.Replace(workId, "/", ":", 1)
	workId = string(unicode.ToUpper(rune(workId[0]))) + workId[1:]
	workId = strings.Replace(workId, "s", "", 1)
	return workId
}
