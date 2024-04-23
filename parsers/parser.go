package parsers

import (
	"net/http"

	"github.com/FlyingButterTuna/wn-tracker/novel"
	"github.com/PuerkitoBio/goquery"
)

type NovelParser interface {
	ParseTOC() ([]novel.SectionData, error)
	ParseTitle() (string, error)
	ParseAuthor() (string, error)
	ParseChapterHtml(chapterPage *goquery.Document) (string, error)
}

type PageFetcher interface {
	FetchPage(url string) (*http.Response, error)
}

const (
	HostNarou    = "ncode.syosetu.com"
	HostKakuyomu = "kakuyomu.jp"
)
