package novel

import "time"

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
	Author   string        `json:"author,omitempty"`
}
