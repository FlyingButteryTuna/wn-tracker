package parsers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/FlyingButterTuna/wn-tracker/novel"
	"github.com/PuerkitoBio/goquery"
)

type NovelFetcher struct {
	client *http.Client
}

func NewFetcher(client *http.Client) *NovelFetcher {
	return &NovelFetcher{client: client}
}

func (p *NovelFetcher) FetchPage(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, nil
}

func (p *NovelFetcher) SaveAllChapters(novelData *novel.NovelData, parser NovelParser, novelDirPath string) error {
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

			resp, err := p.FetchPage(fullChapterLink)
			if err != nil {
				return err
			}

			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				return err
			}

			htmlStr, err := parser.ParseChapterHtml(doc)
			if err != nil {
				return err
			}

			file, err := os.Create(filepath.Join(novelDirPath, chapter.Title+".html"))
			if err != nil {
				return err
			}

			_, err = file.Write([]byte(htmlStr))
			if err != nil {
				return err
			}

			file.Close()
		}
	}
	return nil
}
