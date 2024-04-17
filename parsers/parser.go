package parsers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type NovelParser interface {
	ParseTOC(body *goquery.Document) ([]SectionData, error)
	ParseTitle(body *goquery.Document) (string, error)
}

type ChapterData struct {
	Title       string    `json:"name,omitempty"`
	Link        string    `json:"link,omitempty"`
	DatePosted  time.Time `json:"date_posted,omitempty"`
	DateUpdated time.Time `json:"date_updated,omitempty"`
}

type SectionData struct {
	Name     string        `json:"name,omitempty"`
	Chapters []ChapterData `json:"chapters,omitempty"`
	Level    uint8         `json:"level,omitempty"`
}

type NovelData struct {
	Title    string        `json:"title,omitempty"`
	Sections []SectionData `json:"sections,omitempty"`
	Link     string        `json:"link,omitempty"`
}

const (
	HostNarou    = "ncode.syosetu.com"
	HostKakuyomu = "kakuyomu.jp"
)

func NewParser(urlStr string) (NovelParser, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var parser NovelParser
	switch parsedURL.Host {
	case HostNarou:
		parser = &NarouParser{Link: urlStr}
	case HostKakuyomu:
		parser = &KakuyomuParser{Link: urlStr}
	default:
		return nil, fmt.Errorf("couldn't recognize the host")
	}
	return parser, nil

}

func FetchPage(url string, client *http.Client) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, nil
}

func SaveAllChapters(novelData *NovelData, novelDirPath string, client *http.Client) error {
	_, err := os.Stat(novelDirPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(novelDirPath, 0755)
		if err != nil {
			return err
		}
	}

	for _, section := range novelData.Sections {
		for _, chapter := range section.Chapters {
			fullChapterLink := novelData.Link + chapter.Link

			resp, err := FetchPage(fullChapterLink, client)
			if err != nil {
				return err
			}

			file, err := os.Create(filepath.Join(novelDirPath, chapter.Title+".html"))
			if err != nil {
				return err
			}
			defer file.Close()

			htmlStr, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			file.Write(htmlStr)
		}
	}
	return nil
}
