package kakuyomu

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
	"unicode"

	"github.com/FlyingButterTuna/wn-tracker/novel"
	"github.com/FlyingButterTuna/wn-tracker/parsers"
	"github.com/PuerkitoBio/goquery"
)

type KakuyomuParser struct {
	parsers.CommonParser
	apolloStateJson map[string]interface{}
	link            string
}

const timeLayoutKakuyomu = "2006-01-02T15:04:05Z"

func NewKakuyomuParser(link *url.URL, novelPage *goquery.Document) (*KakuyomuParser, error) {
	apolloStateJson, err := initializeJson(novelPage)
	if err != nil {
		return nil, err
	}
	return &KakuyomuParser{link: link.String(), apolloStateJson: apolloStateJson}, nil
}

func (p *KakuyomuParser) ParseChapterHtml(chapterPage *goquery.Document) (string, error) {
	result := chapterPage.Find("div.widget-episodeBody.js-episode-body[data-viewer-history-path]")
	if result.Length() == 0 {
		return "", fmt.Errorf("chapter text not found")
	}

	p.RemoveAttrsFromElement("p", result)
	p.RemoveRP(result)

	emphasisSpans := result.Find("em.emphasisDots")
	p.ReplaceAttrInElement(emphasisSpans, "class", "em-dot")

	resultStr, err := result.Html()
	if err != nil {
		return "", err
	}

	return strings.Trim(resultStr, "\n"), nil
}

func (p *KakuyomuParser) ParseAuthor() (string, error) {
	novelDataJson, ok := p.apolloStateJson[p.workId()].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("error parsing novel data json")
	}

	authorRef, ok := novelDataJson["author"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("error author id")
	}

	author, ok := p.apolloStateJson[authorRef["__ref"].(string)].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("error author json")
	}

	authorName, ok := author["activityName"].(string)
	if !ok {
		return "", fmt.Errorf("error author name")
	}
	return authorName, nil
}

func (p *KakuyomuParser) ParseTitle() (string, error) {
	novelDataJson, ok := p.apolloStateJson[p.workId()].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("error parsing novel data json")
	}

	novelTitle, ok := novelDataJson["title"].(string)
	if !ok {
		return "", fmt.Errorf("error parsing title from json")
	}
	return novelTitle, nil
}

func (p *KakuyomuParser) ParseTOC() ([]novel.SectionData, error) {
	sections := make([]novel.SectionData, 0)

	toc, ok := p.apolloStateJson[p.workId()].(map[string]interface{})["tableOfContents"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("error parsing novel data json")
	}

	emptyTocStr := "TableOfContentsChapter:"
	var is_sectionless bool
	if firstTocEntry := toc[0].(map[string]interface{})["__ref"].(string); len(toc) == 1 && firstTocEntry == emptyTocStr {
		is_sectionless = true
	}

	// omit error checks for brevity, if toc is found - assume that there are no changes in the json structure
	for _, sectionId := range toc {
		section, _ := p.apolloStateJson[sectionId.(map[string]interface{})["__ref"].(string)].(map[string]interface{})

		chapters := make([]novel.ChapterData, 0)
		chapterIds := section["episodeUnions"].([]interface{})
		for _, chapterId := range chapterIds {
			chapter := chapterId.(map[string]interface{})["__ref"].(string)
			chapterInfo := p.apolloStateJson[chapter].(map[string]interface{})
			chapterTitle := chapterInfo["title"].(string)
			publishedAtStr := chapterInfo["publishedAt"].(string)
			episodeId := chapterInfo["id"].(string)

			chapterLink := "/episodes/" + episodeId
			datePusblished, _ := time.Parse(timeLayoutKakuyomu, publishedAtStr)

			chapterData := novel.ChapterData{Title: chapterTitle, Link: chapterLink, DatePosted: datePusblished}
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

		sectionData := novel.SectionData{Name: sectionTitle, Chapters: chapters, Level: sectionLevel}
		sections = append(sections, sectionData)
	}

	return sections, nil
}

func initializeJson(novelPage *goquery.Document) (map[string]interface{}, error) {
	jsonDataString := novelPage.Find(`[type="application/json"]`).Text()

	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonDataString), &jsonData); err != nil {
		return nil, fmt.Errorf("error parsing JSON data: %w", err)
	}

	props, ok := jsonData["props"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("error parsing 'props' from JSON data")
	}

	pageProps, ok := props["pageProps"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("error parsing 'pageProps' from JSON data")
	}

	apolloState, ok := pageProps["__APOLLO_STATE__"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("error parsing '__APOLLO_STATE__' from JSON data")
	}

	return apolloState, nil
}

func (p *KakuyomuParser) workId() string {
	workId := p.link[strings.LastIndex(p.link[:strings.LastIndex(p.link, "/")-1], "/")+1:]
	workId = strings.Replace(workId, "/", ":", 1)
	workId = string(unicode.ToUpper(rune(workId[0]))) + workId[1:]
	workId = strings.Replace(workId, "s", "", 1)
	return workId
}
