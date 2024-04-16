package parsers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type NovelParser interface {
	ParseTOC(body *goquery.Document) []SectionData
	ParseTitle(body *goquery.Document) string
}

type ChapterData struct {
	Name        string    `json:"name,omitempty"`
	Link        string    `json:"link,omitempty"`
	DatePosted  time.Time `json:"date_posted,omitempty"`
	DateUpdated time.Time `json:"date_updated,omitempty"`
}

type SectionData struct {
	Name     string        `json:"name,omitempty"`
	Chapters []ChapterData `json:"chapters,omitempty"`
}

type NovelData struct {
	Title    string        `json:"title,omitempty"`
	Sections []SectionData `json:"sections,omitempty"`
	Link     string        `json:"link,omitempty"`
}

func FetchPage(url string, client *http.Client) (*goquery.Document, error) {
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

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't load the HTML document")
	}

	return doc, nil
}
