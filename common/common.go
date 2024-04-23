package common

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/FlyingButterTuna/wn-tracker/parsers"
	"github.com/FlyingButterTuna/wn-tracker/parsers/kakuyomu"
	"github.com/FlyingButterTuna/wn-tracker/parsers/narou"
	"github.com/PuerkitoBio/goquery"
)

func NewNovelParser(urlStr string, novelPage *goquery.Document) (parsers.NovelParser, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var parser parsers.NovelParser
	switch parsedURL.Host {
	case parsers.HostNarou:
		parser = narou.NewNarouParser(parsedURL, novelPage, parsers.NewFetcher(&http.Client{}))
	case parsers.HostKakuyomu:
		parser, err = kakuyomu.NewKakuyomuParser(parsedURL, novelPage)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("couldn't recognize the host")
	}
	return parser, nil
}
