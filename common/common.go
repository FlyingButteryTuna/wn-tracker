package common

import (
	"fmt"
	"net/url"

	"github.com/FlyingButterTuna/wn-tracker/parsers"
	"github.com/FlyingButterTuna/wn-tracker/parsers/kakuyomu"
	"github.com/FlyingButterTuna/wn-tracker/parsers/narou"
)

func NewNovelParser(urlStr string) (parsers.NovelParser, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var parser parsers.NovelParser
	switch parsedURL.Host {
	case parsers.HostNarou:
		parser = narou.NewNarouParser(parsedURL)
	case parsers.HostKakuyomu:
		parser = kakuyomu.NewKakuyomuParser(parsedURL)
	default:
		return nil, fmt.Errorf("couldn't recognize the host")
	}
	return parser, nil
}
