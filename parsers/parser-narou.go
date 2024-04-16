package parsers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const timeLayout = "2006/01/02 15:04"

type NarouParser struct {
	Link string
}

func (p *NarouParser) ParseTitle(doc *goquery.Document) string {
	return doc.Find("p.novel_title").Text()
}

func (p *NarouParser) ParseTOC(doc *goquery.Document) []SectionData {
	result := make([]SectionData, 0)

	indexBox := doc.Find(".index_box")

	chapterCounter := 0
	if indexBox.Find(".chapter_title").Length() == 0 {
		result = append(result, SectionData{})
		result[0].Chapters = make([]ChapterData, 0)
		result[0].Name = "default"
		chapterCounter++
	}

	if doc.Find(".novelview_pager").Length() != 0 {
		client := &http.Client{}
		for i := 0; indexBox.Length() != 0; i++ {
			parsePage(&result, &chapterCounter, indexBox)

			doc, err := FetchPage(p.Link+"?p="+strconv.Itoa(i), client)
			if err != nil {
				break
			}
			indexBox = doc.Find(".index_box")
		}
	} else {
		parsePage(&result, &chapterCounter, indexBox)
	}

	return result
}

func parsePage(result *[]SectionData, chapterCounter *int, indexBox *goquery.Selection) {
	indexBox.Children().Each(func(i int, s *goquery.Selection) {
		if s.HasClass("chapter_title") {
			chapterTitle := s.Text()
			(*result) = append((*result), SectionData{})
			(*result)[(*chapterCounter)].Chapters = make([]ChapterData, 0)
			(*result)[(*chapterCounter)].Name = chapterTitle
			(*chapterCounter)++
		} else {
			chapterData := ChapterData{}
			aElem := s.Find(".subtitle").Find("a")
			chapterData.Name = aElem.Text()
			chapterData.Link, _ = aElem.Attr("href")

			updateElem := s.Find(".long_update")
			chapterData.DatePosted, _ = time.Parse(timeLayout, strings.TrimSpace(updateElem.Text()[:17]))
			longUpdate := updateElem.Find("span")
			if longUpdate.Length() != 0 {
				chapterData.DateUpdated, _ = time.Parse(timeLayout, strings.TrimSpace(longUpdate.AttrOr("title", "")[:17]))
			}

			(*result)[(*chapterCounter)-1].Chapters = append((*result)[(*chapterCounter)-1].Chapters, chapterData)
		}
	})
}
