package narou

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/FlyingButterTuna/wn-tracker/novel"
	"github.com/FlyingButterTuna/wn-tracker/parsers"
	"github.com/PuerkitoBio/goquery"
)

const timeLayoutNarou = "2006/01/02 15:04"

type fetcher interface {
	FetchPage(url string) (*http.Response, error)
}

type NarouParser struct {
	fetcher
	parsers.CommonParser
	link string
}

func NewNarouParser(link *url.URL, f fetcher) *NarouParser {
	return &NarouParser{link: link.String(), fetcher: f}
}

func (p *NarouParser) ParseChapterHtml(doc *goquery.Document) (string, error) {
	result := doc.Find("#novel_honbun")
	if result.Length() == 0 {
		return "", fmt.Errorf("chapter text not found")
	}

	p.RemoveAttrsFromElement("p", result)
	p.RemoveRP(result)
	removeRubyWrap(result)

	resultStr, err := result.Html()
	if err != nil {
		return "", err
	}

	return strings.Trim(resultStr, "\n"), nil
}

func (p *NarouParser) ParseAuthor(doc *goquery.Document) (string, error) {
	authorElem := doc.Find(".novel_writername")
	if authorElem.Length() == 0 {
		return "", fmt.Errorf("error parsing the novel_writername element")
	}

	aElem := authorElem.Find("a")
	if aElem.Length() != 0 {
		return aElem.Text(), nil
	}

	authorName := strings.Trim(authorElem.Text(), " \n")
	authorName = authorName[strings.Index(authorName, "：")+3:]
	return authorName, nil
}

func (p *NarouParser) ParseTitle(doc *goquery.Document) (string, error) {
	titleElem := doc.Find("p.novel_title")
	if titleElem.Length() == 0 {
		return "", fmt.Errorf("error parsing the novel_title element")
	}

	return titleElem.Text(), nil
}

func (p *NarouParser) ParseTOC(doc *goquery.Document) ([]novel.SectionData, error) {
	result := make([]novel.SectionData, 0)

	indexBox := doc.Find(".index_box")

	chapterCounter := 0
	if indexBox.Find(".chapter_title").Length() == 0 {
		result = append(result, novel.SectionData{})
		result[0].Chapters = make([]novel.ChapterData, 0)
		result[0].Name = "default"
		chapterCounter++
	}

	if doc.Find(".novelview_pager").Length() != 0 {
		for i := 2; i != 3; i++ {
			parsePage(&result, &chapterCounter, indexBox)

			resp, err := p.fetcher.FetchPage(p.link + "?p=" + strconv.Itoa(i))
			if err != nil {
				break
			}

			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				return nil, err
			}
			indexBox = doc.Find(".index_box")
		}
	} else {
		parsePage(&result, &chapterCounter, indexBox)
	}

	return result, nil
}

func parsePage(result *[]novel.SectionData, chapterCounter *int, indexBox *goquery.Selection) {
	indexBox.Children().Each(func(i int, s *goquery.Selection) {
		if s.HasClass("chapter_title") {
			chapterTitle := s.Text()
			(*result) = append((*result), novel.SectionData{})
			(*result)[(*chapterCounter)].Chapters = make([]novel.ChapterData, 0)
			(*result)[(*chapterCounter)].Name = chapterTitle
			(*chapterCounter)++
		} else {
			chapterData := novel.ChapterData{}

			aElem := s.Find(".subtitle").Find("a")
			chapterData.Title = aElem.Text()

			fullLink, _ := aElem.Attr("href")
			chapterData.Link = fullLink[strings.LastIndex(fullLink[:len(fullLink)-1], "/"):]

			updateElem := s.Find(".long_update")
			chapterData.DatePosted, _ = time.Parse(timeLayoutNarou, strings.TrimSpace(updateElem.Text()[:17]))
			longUpdate := updateElem.Find("span")
			if longUpdate.Length() != 0 {
				chapterData.DateUpdated, _ = time.Parse(timeLayoutNarou, strings.TrimSpace(longUpdate.AttrOr("title", "")[:17]))
			}

			(*result)[(*chapterCounter)-1].Chapters = append((*result)[(*chapterCounter)-1].Chapters, chapterData)
		}
	})
}

func removeRubyWrap(doc *goquery.Selection) {
	doc.Find(".ruby-wrap").Each(func(_ int, s *goquery.Selection) {
		s.Contents().Unwrap()
	})
}
