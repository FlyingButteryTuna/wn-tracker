package parsers

import (
	"github.com/FlyingButterTuna/wn-tracker/novel"
	"github.com/PuerkitoBio/goquery"
)

type NovelParser interface {
	ParseTOC(body *goquery.Document) ([]novel.SectionData, error)
	ParseTitle(body *goquery.Document) (string, error)
}

const (
	HostNarou    = "ncode.syosetu.com"
	HostKakuyomu = "kakuyomu.jp"
)
